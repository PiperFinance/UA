package conf

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

var (
	RunAsServer       bool
	RunAsClient       bool
	RunAsScheduler    bool
	QueueScheduler    *asynq.Scheduler
	QueueClient       *asynq.Client
	QueueServer       *asynq.Server
	asyncQRedisClient asynq.RedisClientOpt
	mux               *asynq.ServeMux
)

//	type queueStatus struct {
//		Client    bool
//		Worker    bool
//		Scheduler bool
//	}
type QueueSchedules struct {
	Cron    string
	Key     string
	Payload []byte
	Q       asynq.Option
	Timeout time.Duration
}

type MuxHandler struct {
	Key     string
	Handler func(context.Context, *asynq.Task) error
	Q       asynq.Option
}

const (
	UsersQ       = "users"
	ScanQ        = "scan"
	FetchQ       = "fetch"
	ParseQ       = "Parse"
	ProcessQ     = "Process"
	UASyncBalQ   = "UA:SyncBal"
	UASyncTHQ    = "UA:SyncTH"
	UASyncNTQ    = "UA:SyncNT"
	MainQ        = "main"
	DefaultQ     = "default"
	UnImportantQ = "Un-Important"
)

func LoadQueue() {
	// Create and configuring Redis connection.
	asyncQRedisClient = asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", Config.RedisUrl.Hostname(), Config.RedisUrl.Port()),
		DB:   Config.RedisDB,
	}
	QueueClient = asynq.NewClient(asyncQRedisClient)

	// Run worker server.
	QueueServer = asynq.NewServer(asyncQRedisClient, asynq.Config{
		Concurrency:  int(Config.MaxConcurrency),
		ErrorHandler: &QueueErrorHandler{},
		Queues: map[string]int{
			UASyncBalQ: 9,
			UASyncNTQ:  5,
			UASyncTHQ:  2,
		},
	})
	mux = asynq.NewServeMux()
	// Block Related

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		Logger.Panic(err)
	}
	QueueScheduler = asynq.NewScheduler(
		asyncQRedisClient,
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)
}

func RunClient() {
	RunAsClient = true
}

func RunWorker(muxHandler []MuxHandler) {
	RunAsServer = true
	for _, mh := range muxHandler {
		mux.HandleFunc(mh.Key, mh.Handler)
	}
	if err := QueueServer.Run(mux); err != nil {
		Logger.Panic(err)
	}
}

func RunScheduler(queueSchedules []QueueSchedules) {
	RunAsScheduler = true
	for _, qs := range queueSchedules {
		_, err := QueueScheduler.Register(qs.Cron, asynq.NewTask(qs.Key, qs.Payload), qs.Q, asynq.Timeout(qs.Timeout))
		if err != nil {
			Logger.Panicf("QueueScheduler: %s", err)
		}
	}
	if err2 := QueueScheduler.Start(); err2 != nil {
		Logger.Panic(err2)
	}
}
