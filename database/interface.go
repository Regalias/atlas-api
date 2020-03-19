package database

import "github.com/regalias/atlas-api/models"

// Provider is the generic interface for interacting with underlying persistent database storage
type Provider interface {

	// InitDatabase is a helper function to initilize the database and/or schema
	InitDatabase() error

	// Getter
	// Returns NotFound error if query return is empty, or operational errors
	GetLinkDetails(linkpath string) (*models.LinkModel, error)

	// TODO: list operation

	// Standard CRUD operations
	// CreateLink creates a new link in the underlying database
	CreateLink(linkmodel *models.LinkModel) error

	// UpdateLink updates the link in the database to match the new model
	// Must return an error if the link does not exist
	UpdateLink(linkmodel *models.LinkModel) error

	// DeleteLink deletes the link from the database
	// Must return an error if the link does not exist
	DeleteLink(linkpath string) error
}
