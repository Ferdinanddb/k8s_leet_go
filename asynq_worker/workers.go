package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"k8s_leet_code_asynq_worker/env"
	"k8s_leet_code_asynq_worker/task"
)

func main() {

	env.LoadEnv()

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	srv := asynq.NewServer(
		asynq.RedisClusterClientOpt{
			Addrs: []string{fmt.Sprintf("%v:%v", redisHost, redisPort)},
			Password: redisPassword,
		},
		asynq.Config{
			Concurrency: 4,
			// Queues: map[string]int{
			// 	"queue:code_request": 4, // processed 100% of the time
			// },
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeRunCode, task.HandleRunCodePythonTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
