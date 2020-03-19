package cache

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis/v7"
)

// RedisProvider contains the context for the cache
type RedisProvider struct {
	client *redis.Client
}

// NewRedisProvider creates a new redis cache provider
func NewRedisProvider(host string, port uint16, opts interface{}) (*RedisProvider, error) {

	// TODO: bring in opts
	r := &RedisProvider{
		client: redis.NewClient(&redis.Options{
			Addr:     host + ":" + strconv.Itoa(int(port)),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	_, err := r.client.Ping().Result()
	// fmt.Printf("%s\n", res)
	return r, err
}

// TODO: handle retry logic?

// FetchLink fetches a linkpath from redis
func (r *RedisProvider) FetchLink(linkpath string) (string, error) {
	val, err := r.client.Get(linkpath).Result()
	if err == redis.Nil {
		// Key does not exist yet
		return "", errors.New("NotFound")
	}
	return val, err
}

// DeleteLink deletes the linkpath key from redis
func (r *RedisProvider) DeleteLink(linkpath string) error {
	_, err := r.client.Del(linkpath).Result()
	// if err == nil && val < 1 {
	// 	return err
	// }
	return err
}

// UpsertLink creates or updates the linkpath key in redis
func (r *RedisProvider) UpsertLink(linkpath string, dest string) error {
	err := r.client.Set(linkpath, dest, 0).Err()
	return err
}
