package main

import (
	"context"
	"time"

	"monitor-service/conf"
	"monitor-service/helpers"
	"monitor-service/internal/server"
)

func main() {
	conf.Initialize()
	srv := server.Init()
	wait := helpers.GracefulShutdown(context.Background(), 10*time.Second, map[string]helpers.Operation{
		"api": func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	<-wait
}
