package redisop

import "github.com/redis/go-redis/v9"

func NewRedisOperator(client redis.Cmdable) Operator {
	return &redisOp{
		cmdable: client,
	}
}

type redisOp struct {
	cmdable redis.Cmdable
}
