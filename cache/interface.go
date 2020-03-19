package cache

// Provider is the generic interface for interacting with an underlying cache provider
type Provider interface {

	// FetchLink attempts to grab a link key from the cache
	// Returns a NotFound error if the key is empty or does not exist
	FetchLink(linkpath string) (string, error)

	// DeleteLink will remove the linkpath key from the cache
	// Returns an error only on operational errors
	DeleteLink(linkpath string) error

	// UpsertLink will insert a linkpath:dest mapping into the cache
	// Returns an error only on operational errors
	UpsertLink(linkpath string, dest string) error
}
