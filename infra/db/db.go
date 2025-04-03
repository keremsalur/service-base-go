package db

import "gorm.io/gorm"

type Database interface {
	GetDB() *gorm.DB
	Connect(dsn string) Database
	Migrate(models ...interface{})
}
