package cache

import (
	"errors"
	"time"

	"github.com/allegro/bigcache"
)

// LocalProvider contains context for a local in-memory cache
type LocalProvider struct {
	cache *bigcache.BigCache
}

// NewLocalProvider creates a new local in-memory cache
func NewLocalProvider(expiryTimeSec int64) (*LocalProvider, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(expiryTimeSec) * time.Second))

	lp := &LocalProvider{
		cache: cache,
	}

	return lp, err
}

// FetchLink attempts to grab a link key from the cache
// Returns a NotFound error if the key is empty or does not exist
func (lp *LocalProvider) FetchLink(linkpath string) (string, error) {

	entry, err := lp.cache.Get(linkpath)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return "", errors.New("NotFound")
		}
		return "", err
	}
	// Additional check that key is not empty
	dest := string(entry)
	if len(dest) < 1 {
		return "", errors.New("NotFound")
	}
	return dest, err
}

// DeleteLink will remove the linkpath key from the cache
// Returns an error only on operational errors
func (lp *LocalProvider) DeleteLink(linkpath string) error {
	return lp.cache.Delete(linkpath)
}

// UpsertLink will insert a linkpath:dest mapping into the cache
// Returns an error only on operational errors
func (lp *LocalProvider) UpsertLink(linkpath string, dest string) error {
	return lp.cache.Set(linkpath, []byte(dest))
}
