package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB(dsn string) (*DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
