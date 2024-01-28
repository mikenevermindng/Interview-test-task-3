//go:build wireinject
// +build wireinject

package schedule

import (
	"github.com/google/wire"

	"monitor-service/conf"
	"monitor-service/internal/db"
)

var ProviderScheduleSet = wire.NewSet(conf.NewConfiguration, db.NewDatabase, NewSchedule)

func CreateSchedule() (*Schedule, error) {
	wire.Build(ProviderScheduleSet)

	return &Schedule{}, nil
}
