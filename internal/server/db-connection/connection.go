package db_connection

import (
	"gorm.io/gorm"

	"monitor-service/internal/db"
)

var database = (*db.Database)(nil)

func ConnectDatabase() {
	newDb, _ := db.CreateDatabaseClient()
	database = newDb
}

func GetDB() *gorm.DB {
	return database.Client
}
