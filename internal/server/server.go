package server

import (
	api_conf "monitor-service/internal/server/api-conf"
	dbconnection "monitor-service/internal/server/db-connection"
)

func Init() {
	r := NewRouter()

	dbconnection.ConnectDatabase()
	api_conf.LoadConfig()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
