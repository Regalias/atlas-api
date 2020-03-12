package apiserver

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func newAWSSession() *session.Session {
	sess := session.Must(session.NewSession())
	return sess
}

// ensureTable attempts to describe the requested table, and creates one if it doesn't exist
func (dp *DataProvider) ensureTable(tableName string) error {
	_, err := dp.ddb.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				dp.logger.Debug().Msg(dynamodb.ErrCodeResourceNotFoundException + ":" + aerr.Error())
				// Table doesn't exist, lets create it
				dp.logger.Info().Msg("Table " + tableName + " not found, creating it now...")
				if err := dp.createTable(tableName); err != nil {
					return err
				}

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
	}
	return err
}

// createTable creates the target DDB table with the required schema
func (dp *DataProvider) createTable(tableName string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Artist"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SongTitle"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Artist"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SongTitle"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	_, err := dp.ddb.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				dp.logger.Error().Msg(dynamodb.ErrCodeResourceInUseException + ":" + aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				dp.logger.Error().Msg(dynamodb.ErrCodeLimitExceededException + ":" + aerr.Error())
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
	}
	return err
}
