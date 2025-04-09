package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {

	db, err := gorm.Open(postgres.Open(DB_URL), &gorm.Config{})
	// log.Println("Connected to database",DB_URL)
	if err != nil {
		log.Println("Found err while connecting to database", err)
	}

	return db
}

func CloseDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Println("Found err while closing the database", err)
	}
	dbSQL.Close()
}
