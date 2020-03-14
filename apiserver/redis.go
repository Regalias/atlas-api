package apiserver

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis/v7"
)

type cacheProvider struct {
	client *redis.Client
}

func newCache(host string, port uint16, opts interface{}) (*cacheProvider, error) {
	c := &cacheProvider{
		client: redis.NewClient(&redis.Options{
			Addr:     host + ":" + strconv.Itoa(int(port)),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	_, err := c.client.Ping().Result()
	// fmt.Printf("%s\n", res)
	return c, err
}

// TODO: handle retry logic?
func (c *cacheProvider) fetchLink(linkpath string) (string, error) {
	val, err := c.client.Get(linkpath).Result()
	if err == redis.Nil {
		// Key does not exist yet
		return "", errors.New("NotFound")
	}
	return val, err
}

func (c *cacheProvider) deleteLink(linkpath string) error {
	_, err := c.client.Del(linkpath).Result()
	// if err == nil && val < 1 {
	// 	return err
	// }
	return err
}

func (c *cacheProvider) upsertLink(linkpath string, dest string) error {
	err := c.client.Set(linkpath, dest, 0).Err()
	return err
}
