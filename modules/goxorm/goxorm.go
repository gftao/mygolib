package goxorm

import (
	"mygolib/modules/config"
	"errors"
)

var instance *Enginesql

func InitModel() error {
	var err error
	instance, err = initDb()
	return err

}

func initDb() (*Enginesql, error) {
	dbconf := config.StringDefault("dbconf", "")
	if dbconf == "" {
		return nil, errors.New("数据库配置文件为空"+dbconf)
	}

	db, err := opendb(dbconf)
	if err != nil {
		return nil, err
	}
	return db, err
}
