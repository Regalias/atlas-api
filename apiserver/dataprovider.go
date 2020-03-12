package apiserver

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
)

// DataProvider contains methods to interact with the underlying database used for persistent storage
type DataProvider struct {
	sess   *session.Session
	ddb    *dynamodb.DynamoDB
	logger *zerolog.Logger
}

// NewDataProvider creates and configures a new DataProvider object
func NewDataProvider(logger *zerolog.Logger) (*DataProvider, error) {
	sess := newAWSSession()
	ddb := dynamodb.New(sess)

	dp := &DataProvider{
		sess:   sess,
		ddb:    ddb,
		logger: logger,
	}
	return dp, nil
}

// GetLinkDetails fetches the link details based on a link ID
func (dp *DataProvider) GetLinkDetails(linkid string) (*LinkModel, error) {
	return &LinkModel{}, nil
}

// CreateLink creates a new link from the supplied model
func (dp *DataProvider) CreateLink(linkmodel *LinkModel) error {
	return nil
}

// DeleteLink deletes the link matching the link ID in the supplied model
func (dp *DataProvider) DeleteLink(linkmodel *LinkModel) error {
	return nil
}

// UpdateLink updates the existing link matching the link ID in the supplied model
func (dp *DataProvider) UpdateLink(linkmodel *LinkModel) error {
	return nil
}
