package asynq_client

import (
	"fmt"
	"os"

	// "github.com/redis/go-redis/v9"

	"github.com/hibiken/asynq"
)

// var RedisClient *redis.ClusterClient
var AsynqRedisClient *asynq.Client

func Connect() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{fmt.Sprintf("%v:%v", redisHost, redisPort)},
	// 	// Username: ,
	// 	Password: redisPassword,
	// })

	AsynqRedisClient = asynq.NewClient(asynq.RedisClusterClientOpt{
		Addrs: []string{fmt.Sprintf("%v:%v", redisHost, redisPort)},
		// Username: ,
		Password: redisPassword,
	})
	// defer AsynqRedisClient.Close()
}