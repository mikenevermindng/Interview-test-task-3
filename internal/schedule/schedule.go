package schedule

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-co-op/gocron"

	"monitor-service/conf"
	"monitor-service/internal/db"
	"monitor-service/internal/db/models"
	uribuilder "monitor-service/internal/uri-builder"
	"monitor-service/internal/worker"
)

type Schedule struct {
	quit                  chan bool
	queue                 chan uribuilder.UrlBuilder
	siteUpdateCh          chan uribuilder.UrlBuilder
	sites                 map[string]uribuilder.UrlBuilder
	scheduler             *gocron.Scheduler
	intervalMs            int
	maxMonitorConcurrency int
	maxRedirect           int
	reqTimeout            int
	mu                    sync.Mutex
	db                    *db.Database
}

func NewSchedule(configuration *conf.Configuration, db *db.Database) *Schedule {
	scheduler := gocron.NewScheduler(time.UTC)
	sch := &Schedule{
		quit:                  make(chan bool, configuration.Monitor.MaxMonitorConcurrency),
		queue:                 make(chan uribuilder.UrlBuilder, configuration.Monitor.MaxMonitorConcurrency),
		siteUpdateCh:          make(chan uribuilder.UrlBuilder),
		sites:                 map[string]uribuilder.UrlBuilder{},
		scheduler:             scheduler,
		intervalMs:            configuration.Monitor.MonitorIntervalMs,
		maxMonitorConcurrency: configuration.Monitor.MaxMonitorConcurrency,
		maxRedirect:           configuration.Monitor.MaxRedirect,
		reqTimeout:            configuration.Monitor.RequestTimeout,
		db:                    db,
	}
	return sch
}

func (sch *Schedule) Initial() {
	sch.readListServices()
	sch.setupWorker()
}

func (sch *Schedule) Start() error {

	if _, err := sch.scheduler.Every(sch.intervalMs).Milliseconds().Do(func() {
		for _, site := range sch.sites {
			sch.queue <- site
		}
	}); err != nil {
		fmt.Printf("Unable to start consumer: %v", err)
		panic(err)
	}

	go sch.eventListener()

	sch.scheduler.StartAsync()
	return nil
}

func (sch *Schedule) Stop(ctx context.Context) error {
	close(sch.queue)
	for i := 1; i <= sch.maxMonitorConcurrency; i++ {
		sch.quit <- true
	}
	return nil
}

func (sch *Schedule) eventListener() {
	for {
		select {
		case siteUpdated := <-sch.siteUpdateCh:
			sch.mu.Lock()
			sch.sites[siteUpdated.GetUri()] = siteUpdated
			sch.mu.Unlock()
		}
	}
}

func (sch *Schedule) readListServices() {
	dbClient := sch.db.Client
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get current path")
		panic(err)
	}
	f, err := os.Open(fmt.Sprintf("%s/%s", pwd, "sites.txt"))
	if err != nil {
		fmt.Printf("Unable to read sites.txt: %s\n", err)
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		uri := scanner.Text()
		ssl := true
		var lastHeartbeat models.HeartBeat
		err = dbClient.Model(&models.HeartBeat{}).Where("uri = ?", uri).Order("created_at DESC").First(&lastHeartbeat).Error
		if err == nil {
			ssl = lastHeartbeat.Ssl
		}
		site := uribuilder.UrlBuilder{}
		site.Ssl(ssl)
		site.Uri(uri)
		ok, _ := site.BuildURI()

		if ok {
			sch.sites[uri] = site
		} else {
			fmt.Println("Invalid uri")
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Fail to read content sites.txt: %s\n", err)
		panic(err)
	}
}

func (sch *Schedule) setupWorker() {
	for i := 1; i <= sch.maxMonitorConcurrency; i++ {
		w := worker.NewWorker(i, sch.db, sch.queue, sch.quit, sch.siteUpdateCh, sch.maxRedirect, sch.reqTimeout)
		go w.Start()
	}
}
