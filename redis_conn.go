package main

import "github.com/go-redis/redis"

const (
	rdsAddr = "192.168.174.149:6379"
)

var (
	rdsClient *redis.Client
)

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     rdsAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// Output: PONG <nil>
	return client
}
