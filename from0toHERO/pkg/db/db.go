package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func New() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=steven dbname=fullstack_api password=password")

	if err != nil {
		fmt.Println("connection err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}
