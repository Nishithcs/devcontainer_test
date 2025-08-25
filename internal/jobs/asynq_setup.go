package jobs

import (
	"fmt"
	"github.com/hibiken/asynq"
	"os"
)

func NewAsynqClient() *asynq.Client {
	opt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	return asynq.NewClient(opt)
}

func NewAsynqServer() *asynq.Server {
	opt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	return asynq.NewServer(opt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"default":  6,
			"critical": 10,
		},
	})
}
