package worker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"monitor-service/internal/db"
	"monitor-service/internal/db/models"
	uribuilder "monitor-service/internal/uri-builder"
	myhttp "monitor-service/pkg/my-http"
)

type Worker struct {
	id           int
	queue        <-chan uribuilder.UrlBuilder
	siteUpdateCh chan uribuilder.UrlBuilder
	quit         <-chan bool
	http         *myhttp.MyHttp
	timeoutMs    int
	db           *db.Database
}

func NewWorker(id int, db *db.Database, queue <-chan uribuilder.UrlBuilder, quit <-chan bool, siteUpdateCh chan uribuilder.UrlBuilder, maxRedirect, requestTimeout int) *Worker {
	httpClient := myhttp.New(maxRedirect)

	w := &Worker{
		id:           id,
		queue:        queue,
		quit:         quit,
		http:         httpClient,
		timeoutMs:    requestTimeout,
		siteUpdateCh: siteUpdateCh,
		db:           db,
	}
	return w
}

func (w *Worker) Start() {
	for {
		select {
		case <-w.quit:
			return
		case site := <-w.queue:
			ok, uri := site.BuildURI()
			if !ok {
				fmt.Printf("Invalid site")
			}
			err, responseTime := w.monitor(uri)
			errMessage := (*string)(nil)
			status := "UP"
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					fmt.Printf("Request timeout on %s: %s\n", uri, err)
				case errors.Is(err, http.ErrSchemeMismatch):
					fmt.Printf("Request fail with error scheme mismatch on %s: %s\n", uri, err)
					updateSite := uribuilder.UrlBuilder{}
					updateSite.Uri(site.GetUri())
					updateSite.Ssl(false)
					w.siteUpdateCh <- updateSite
				default:
					fmt.Printf("Error making request on %s: %s\n", uri, err)
				}
				errStr := err.Error()
				errMessage = &errStr
				status = "DOWN"
			}
			fmt.Printf("response time of %s: %s\n", uri, responseTime)

			dbClient := w.db.Client
			dbClient.Model(&models.HeartBeat{}).Create(&models.HeartBeat{
				ResponseTime: int(responseTime.Milliseconds()),
				Uri:          site.GetUri(),
				Ssl:          site.GetSsl(),
				Timeout:      errors.Is(err, context.DeadlineExceeded),
				Error:        errMessage,
				Status:       status,
			})
		}
	}
}

func (w *Worker) monitor(uri string) (error, time.Duration) {
	client := w.http.Client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	start := time.Now()
	resp, err := client.Do(req)
	responseTime := time.Since(start)
	if err != nil {
		return err, responseTime
	}
	defer resp.Body.Close()
	return nil, responseTime
}
