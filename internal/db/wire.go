//go:build wireinject
// +build wireinject

package db

import (
	"github.com/google/wire"

	"monitor-service/conf"
)

var ProviderDBSet = wire.NewSet(conf.NewConfiguration, NewDatabase)

func CreateDatabaseClient() (*Database, error) {
	wire.Build(ProviderDBSet)

	return &Database{}, nil
}
