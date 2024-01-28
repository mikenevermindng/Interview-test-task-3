package main

import (
	"context"
	"time"

	"monitor-service/conf"
	"monitor-service/helpers"
	"monitor-service/internal/schedule"
)

func main() {
	conf.Initialize()
	sch, err := schedule.CreateSchedule()
	if err != nil {
		panic(err)
	}
	sch.Initial()
	sch.Start()

	wait := helpers.GracefulShutdown(context.Background(), 10*time.Second, map[string]helpers.Operation{
		"scheduler": func(ctx context.Context) error {
			return sch.Stop(ctx)
		},
	})

	<-wait
}
