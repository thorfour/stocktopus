package cfg

import "os"

var (
	// RedisAddr is a redis endpoint to store information
	RedisAddr string
	// RedisPw is if the redis endpoint has a password using the AUTH command
	RedisPw string
)

func init() {
	RedisAddr = os.Getenv("REDISADDR")
	RedisPw = os.Getenv("REDISPW")
}
