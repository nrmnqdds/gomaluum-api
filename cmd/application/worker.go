package application

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/nrmnqdds/gomaluum-api/tasks"
)

// StartAsynqServer starts the Asynq worker server.
// This function is called when the -w or --worker flag is provided.
func StartAsynqServer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     helpers.GetEnv("REDIS_URL"),
			Username: helpers.GetEnv("REDIS_USERNAME"),
			Password: helpers.GetEnv("REDIS_PASSWORD"),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeRedisSave, tasks.SaveToRedis)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
	log.Println("Asynq Server started")
}
