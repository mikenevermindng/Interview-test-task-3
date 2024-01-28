package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"monitor-service/conf"
	"monitor-service/internal/db/models"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(config *conf.Configuration) *Database {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open(config.Database.Sqlite), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Printf("Unable to connect database: %v\n", err)
		panic(err)
	}

	err = db.AutoMigrate(
		&models.HeartBeat{},
		&models.UserRequest{},
	)
	if err != nil {
		fmt.Printf("Unable to migrate database: %v\n", err)
		panic(err)
	}

	database := &Database{
		Client: db,
	}
	return database
}
