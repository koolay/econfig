// Package store provides ...
package store

import "github.com/koolay/econfig/context"

// import "github.com/spf13/viper"

type Storage interface {
	GetItem(key string) (string, error)
	SetItem(key string, value interface{}) error
	GetItems(keys []string) (map[string]interface{}, error)
}

func NewStorage() Storage {
	context.Logger.INFO.Println("use redis store")
	return newRedisStorage("localhost:6379", "", 0)
}
