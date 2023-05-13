package redisop

import "github.com/redis/go-redis/v9"

func NewRedisOperator(client *redis.Client) Operator {
	return &redisOp{
		client: client,
	}
}

type redisOp struct {
	client *redis.Client
}
