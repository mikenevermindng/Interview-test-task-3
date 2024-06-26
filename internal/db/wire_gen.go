// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package db

import (
	"github.com/google/wire"
	"monitor-service/conf"
)

// Injectors from wire.go:

func CreateDatabaseClient() (*Database, error) {
	configuration := conf.NewConfiguration()
	database := NewDatabase(configuration)
	return database, nil
}

// wire.go:

var ProviderDBSet = wire.NewSet(conf.NewConfiguration, NewDatabase)
