package redisop

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	sdTest "github.com/toysmars/distlib-go/pkg/servicediscovery/test"
	queueTest "github.com/toysmars/distlib-go/pkg/task/queue/test"
)

func createTestRedisClient() (*redis.Client, error) {
	s, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})
	return client, nil
}

func TestRedisQueueOperator(t *testing.T) {
	t.Parallel()

	client, err := createTestRedisClient()
	assert.NoError(t, err)

	queueTest.RunOperatorUnitTests(t, NewRedisOperator(client))
}

func TestRedisServiceDiscoveryOperator(t *testing.T) {
	t.Parallel()

	client, err := createTestRedisClient()
	assert.NoError(t, err)

	sdTest.RunOperatorUnitTests(t, NewRedisOperator(client))
}
