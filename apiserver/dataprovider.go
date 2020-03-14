package apiserver

import (
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog"
)

// DataProvider contains methods to interact with the underlying database used for persistent storage
type DataProvider struct {
	sess      *session.Session
	ddb       *dynamodb.DynamoDB
	logger    *zerolog.Logger
	tableName string
}

// NewDataProvider creates and configures a new DataProvider object
func NewDataProvider(logger *zerolog.Logger, tableName string) (*DataProvider, error) {
	sess := newAWSSession()
	dp := &DataProvider{
		sess:      sess,
		ddb:       dynamodb.New(sess),
		logger:    logger,
		tableName: tableName,
	}
	return dp, nil
}

// GetLinkDetails fetches the link details based on a link path
func (dp *DataProvider) GetLinkDetails(linkpath string) (*LinkModel, error) {
	resp, err := dp.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dp.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"LinkPath": {
				S: aws.String(linkpath),
			},
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeProvisionedThroughputExceededException + ":" + aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				dp.logger.Error().Msg(dynamodb.ErrCodeResourceNotFoundException + ":" + aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				dp.logger.Error().Msg(dynamodb.ErrCodeRequestLimitExceeded + ":" + aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				dp.logger.Error().Msg(dynamodb.ErrCodeInternalServerError + ":" + aerr.Error())
			default:
				dp.logger.Error().Msg(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			dp.logger.Error().Msg(err.Error())
		}
		return nil, err
	}

	lm := &LinkModel{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &lm)
	if err != nil {
		dp.logger.Error().Msg("Failed to unmarshal Record: " + err.Error())
		return nil, err
	}

	if lm.LinkPath == "" {
		return nil, errors.New("NotFound")
	}

	return lm, nil
}

// CreateLink creates a new link from the supplied model
func (dp *DataProvider) CreateLink(linkmodel *LinkModel) error {
	link, err := dynamodbattribute.MarshalMap(*linkmodel)

	if err != nil {
		dp.logger.Error().Msg("DDB Marshal Failed: " + err.Error())
		return err
	}

	_, err = dp.ddb.PutItem(&dynamodb.PutItemInput{
		Item:                link,
		TableName:           aws.String(dp.tableName),
		ConditionExpression: aws.String("attribute_not_exists(LinkPath)"), // must be unique
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				dp.logger.Error().Msg(dynamodb.ErrCodeConditionalCheckFailedException + ":" + aerr.Error())
				// Not unique
				return errors.New("AlreadyExists")
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeProvisionedThroughputExceededException + ":" + aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				dp.logger.Error().Msg(dynamodb.ErrCodeResourceNotFoundException + ":" + aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeItemCollectionSizeLimitExceededException + ":" + aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				dp.logger.Error().Msg(dynamodb.ErrCodeTransactionConflictException + ":" + aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				dp.logger.Error().Msg(dynamodb.ErrCodeRequestLimitExceeded + ":" + aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				dp.logger.Error().Msg(dynamodb.ErrCodeInternalServerError + ":" + aerr.Error())
			default:
				dp.logger.Error().Msg("DDB PutItem Failed: " + aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			dp.logger.Error().Msg("DDB PutItem Failed: " + err.Error())
		}
	}
	return err
}

// DeleteLink deletes the link matching the link path in the supplied model
func (dp *DataProvider) DeleteLink(linkpath string) error {
	// DeleteItem is idempotent - need to specify a condition that it must exist to be successful
	_, err := dp.ddb.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(dp.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"LinkPath": {
				S: aws.String(linkpath),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":lp": {
				S: aws.String(linkpath),
			},
		},
		ConditionExpression: aws.String("LinkPath = :lp"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.New("NotFound")
			default:
				dp.logger.Error().Msg("DDB DeleteItem Failed: " + aerr.Error())
				return err
			}
		} else {
			dp.logger.Error().Msg("DDB DeleteItem Failed: " + err.Error())
		}
	}
	return err
}

// UpdateLink updates the existing link matching the link path in the supplied model
func (dp *DataProvider) UpdateLink(linkmodel *LinkModel) error {

	// Query the existing link to check for existance and differences
	res, err := dp.GetLinkDetails(linkmodel.LinkPath)
	if err != nil {
		return err // pass back upstream error
	}

	if checkLinkModelsAreEqual(linkmodel, res) {
		// Models are same, no changes required!
		return errors.New("NoChange")
	}

	link := map[string]*dynamodb.AttributeValue{
		"LinkPath": {S: aws.String(linkmodel.LinkPath)},
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#CN":  aws.String("CanonicalName"),
			"#TU":  aws.String("TargetURL"),
			"#EN":  aws.String("Enabled"),
			"#LM":  aws.String("LastModified"),
			"#LMB": aws.String("LastModifiedBy"),
		},
		TableName:        aws.String(dp.tableName),
		ReturnValues:     aws.String("NONE"),
		UpdateExpression: aws.String("set #CN = :cn, #TU = :tu, #EN = :en, #LM = :lm, #LMB = :lmb"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":lp": {
				S: aws.String(linkmodel.LinkPath),
			},
			":tu": {
				S: aws.String(linkmodel.TargetURL),
			},
			":cn": {
				S: aws.String(linkmodel.CanonicalName),
			},
			":lm": {
				N: aws.String(strconv.FormatInt(linkmodel.LastModified, 10)),
			},
			":lmb": {
				S: aws.String(linkmodel.LastModifiedBy),
			},
			":en": {
				BOOL: aws.Bool(linkmodel.Enabled),
			},
		},
		ConditionExpression: aws.String("LinkPath = :lp"), // TODO: this does not quite work
		Key:                 link,
	}

	_, err = dp.ddb.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				dp.logger.Debug().Msg(dynamodb.ErrCodeConditionalCheckFailedException + ":" + aerr.Error())
				// Item does not exist - we shouldn't get here as we already checked this before
				return errors.New("NotFound")
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeProvisionedThroughputExceededException + ":" + aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				dp.logger.Error().Msg(dynamodb.ErrCodeResourceNotFoundException + ":" + aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeItemCollectionSizeLimitExceededException + ":" + aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				dp.logger.Error().Msg(dynamodb.ErrCodeTransactionConflictException + ":" + aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				dp.logger.Error().Msg(dynamodb.ErrCodeRequestLimitExceeded + ":" + aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				dp.logger.Error().Msg(dynamodb.ErrCodeInternalServerError + ":" + aerr.Error())
			default:
				dp.logger.Error().Msg(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			dp.logger.Error().Msg(err.Error())
		}
		return err
	}
	return err
}
