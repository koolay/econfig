// Package store provides ...
package store

import (
	"errors"
	"fmt"

	"github.com/koolay/econfig/config"
	"github.com/koolay/econfig/context"
)

// import "github.com/spf13/viper"

type Storage interface {
	GetItem(key string) (string, error)
	SetItem(key string, value interface{}) error
	GetItems(keys []string) (map[string]interface{}, error)
}

func NewStorage(backend string) (Storage, error) {

	context.Logger.INFO.Println("use backend ", backend)
	optionsMap, err := config.GetBackends(backend)
	if err != nil {
		return nil, err
	}

	switch backend {
	case "redis":
		host := config.ValueOfMap("host", optionsMap, "localhost")
		port := config.ValueOfMap("port", optionsMap, "6379")
		pwd := config.ValueOfMap("password", optionsMap, "")
		return newRedisStorage(fmt.Sprintf("%s:%s", host, port), pwd, 0)
	case "mysql", "postgres":
		dsn := config.ValueOfMap("dsn", optionsMap, "")
		return newRSqlStorage(backend, dsn)
	}
	return nil, errors.New("not supported backend")
}
