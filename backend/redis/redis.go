package redis

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.ClusterClient

func Connect() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{fmt.Sprintf("%v:%v", redisHost, redisPort)},
		// Username: ,
		Password: redisPassword,
	})
}