// Package store provides ...
package store

import (
	"github.com/koolay/econfig/context"
	redis "gopkg.in/redis.v5"
)

type redisStorage struct {
	client *redis.Client

	// host:port address.eg: localhost:6379
	Addr string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string

	// Database to be selected after connecting to the server.
	DB int
}

func newRedisStorage(addr string, password string, db int) (*redisStorage, error) {
	return &redisStorage{Addr: addr, Password: password, DB: db}, nil
}

func (rs *redisStorage) connect() *redis.Client {
	redisOpts := &redis.Options{Addr: rs.Addr, Password: rs.Password, DB: rs.DB}
	client := redis.NewClient(redisOpts)
	_, err := client.Ping().Result()

	if err != nil {
		context.Logger.FATAL.Fatalf("cannot connect to redis, %v", err)
	}

	return client
}

func (rs *redisStorage) GetItem(key string) (string, error) {
	client := rs.connect()
	defer client.Close()
	return client.Get(key).Result()
}

func (rs *redisStorage) SetItem(key string, value interface{}) error {
	client := rs.connect()
	defer client.Close()
	return client.Set(key, value, 0).Err()
}

func (rs *redisStorage) GetItems(keys []string) (map[string]interface{}, error) {
	client := rs.connect()
	defer client.Close()
	kv := make(map[string]interface{})
	for _, key := range keys {
		if value, err := client.Get(key).Result(); err == nil {
			context.Logger.INFO.Printf("value: %s", value)
			kv[key] = value
		} else if err != redis.Nil {
			return kv, err
		}
	}
	return kv, nil
}
