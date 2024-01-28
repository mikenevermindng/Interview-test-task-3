package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	api_conf "monitor-service/internal/server/api-conf"
	dbconnection "monitor-service/internal/server/db-connection"
)

func Init() *http.Server {
	r := NewRouter()

	dbconnection.ConnectDatabase()
	api_conf.LoadConfig()
	apiConfig := api_conf.GetConfig()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", apiConfig.Port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}
