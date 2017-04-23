package store

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/koolay/econfig/context"
	_ "github.com/lib/pq"
)

var (
	tableName = "econfig"
)

type rsqlStorage struct {
	driver string
	dsn    string
}

type keyValuePair struct {
	Key   string
	Value string
}

func newRSqlStorage(driver string, dsn string) (*rsqlStorage, error) {
	if dsn == "" {
		return nil, errors.New("dsn not allow empty")
	}
	rsql := &rsqlStorage{dsn: dsn, driver: driver}
	context.Logger.INFO.Println("test sql connection")
	if sess, err := rsql.open(); err == nil {
		defer sess.Close()
		if err = sess.Ping(); err != nil {
			return nil, err
		}
		return rsql, nil
	} else {
		return nil, err
	}
}

func (rs *rsqlStorage) open() (*dbr.Session, error) {
	if conn, err := dbr.Open(rs.driver, rs.dsn, nil); err == nil {
		return conn.NewSession(nil), nil
	} else {
		return nil, err
	}
}

func (rs *rsqlStorage) GetItem(key string) (string, error) {
	var value string
	var err error
	conn, err := rs.open()
	if err == nil {
		defer conn.Close()
		err = conn.Select("value").From(tableName).Where("key=?", key).LoadValue(&value)
	}
	return value, err
}
func (rs *rsqlStorage) SetItem(key string, value interface{}) error {
	var exits int
	conn, err := rs.open()

	if err == nil {
		defer conn.Close()
		if err := conn.Select("count(*)").From(tableName).Where("key=?", key).LoadValue(&exits); err == nil {
			if exits > 0 {
				sqlResult, err := conn.Update(tableName).Where("key=?", key).Set("value", value.(string)).Exec()
				if err == nil {
					affected, err := sqlResult.RowsAffected()
					if err == nil {
						if affected == 0 {
							return errors.New("no affected rows")
						}
					}
				}
				return err
			} else {
				_, err = conn.InsertInto(tableName).Columns("key", "value").Values(key, value).Exec()
			}
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func (rs *rsqlStorage) GetItems(keys []string) (map[string]interface{}, error) {

	keyValuePairs := make(map[string]interface{})
	rows := []keyValuePair{}
	conn, err := rs.open()

	if err == nil {
		defer conn.Close()
		if _, err := conn.Select("key", "value").From(tableName).Where("key in ?", keys).Load(&rows); err == nil {
			for _, item := range rows {
				keyValuePairs[item.Key] = item.Value
			}
		}

	}

	return keyValuePairs, nil
}
