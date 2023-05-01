package queue

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/toysmars/distlib-go/pkg/task"
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

func TestRedisOperator_FIFO(t *testing.T) {
	t.Parallel()

	option := Option{
		Kind: Fifo,
		TaskGroup: task.Group{
			Namespace: "unit-test",
			Name:      "test-redis-op-fifo",
		},
	}
	client, err := createTestRedisClient()
	assert.NoError(t, err)

	ts := createOperatorTestSetup(NewRedisOperator(client, option), option)

	ts.runTestActions(t, []testAction{
		{Action: "pop", ExpectedError: ErrNotFound},
		{Action: "push", Task: ts.createTask("task-1", task.Option{})},
		{Action: "push", Task: ts.createTask("task-2", task.Option{})},
		{Action: "pop", ExpectedMessage: "task-1"},
		{Action: "push", Task: ts.createTask("task-3", task.Option{})},
		{Action: "pop", ExpectedMessage: "task-2"},
		{Action: "pop", ExpectedMessage: "task-3"},
		{Action: "pop", ExpectedError: ErrNotFound},
		{Action: "push", Task: ts.createTask("task-4", task.Option{})},
		{Action: "push", Task: ts.createTask("task-5", task.Option{})},
		{Action: "pop", ExpectedMessage: "task-4"},
		{Action: "push", Task: ts.createTask("task-6", task.Option{})},
		{Action: "pop", ExpectedMessage: "task-5"},
		{Action: "pop", ExpectedMessage: "task-6"},
		{Action: "pop", ExpectedError: ErrNotFound},
	})
}

func TestRedisOperator_Priority(t *testing.T) {
	t.Parallel()

	option := Option{
		Kind: Priority,
		TaskGroup: task.Group{
			Namespace: "unit-test",
			Name:      "test-redis-op-priority",
		},
	}
	client, err := createTestRedisClient()
	assert.NoError(t, err)

	ts := createOperatorTestSetup(NewRedisOperator(client, option), option)

	ts.runTestActions(t, []testAction{
		{Action: "pop", ExpectedError: ErrNotFound},
		{Action: "push", Task: ts.createTask("task-1", task.Option{Priority: 3})},
		{Action: "push", Task: ts.createTask("task-2", task.Option{Priority: 5})},
		{Action: "push", Task: ts.createTask("task-3", task.Option{Priority: 2})},
		{Action: "pop", ExpectedMessage: "task-3"},
		{Action: "push", Task: ts.createTask("task-4", task.Option{Priority: 4})},
		{Action: "pop", ExpectedMessage: "task-1"},
		{Action: "push", Task: ts.createTask("task-5", task.Option{Priority: 1})},
		{Action: "pop", ExpectedMessage: "task-5"},
		{Action: "pop", ExpectedMessage: "task-4"},
		{Action: "pop", ExpectedMessage: "task-2"},
		{Action: "pop", ExpectedError: ErrNotFound},
	})
}

func TestRedisOperator_Scheduled(t *testing.T) {
	t.Parallel()

	option := Option{
		Kind: Scheduled,
		TaskGroup: task.Group{
			Namespace: "unit-test",
			Name:      "test-redis-op-scheduled",
		},
	}
	client, err := createTestRedisClient()
	assert.NoError(t, err)

	now := time.Now()
	ts := createOperatorTestSetup(NewRedisOperator(client, option), option)

	ts.runTestActions(t, []testAction{
		{Action: "pop", ExpectedError: ErrNotFound},
		{Action: "push", Task: ts.createTask("task-1", task.Option{ScheduledFor: now.Add(300 * time.Millisecond)})},
		{Action: "push", Task: ts.createTask("task-2", task.Option{ScheduledFor: now.Add(200 * time.Millisecond)})},
		{Action: "push", Task: ts.createTask("task-3", task.Option{ScheduledFor: now.Add(100 * time.Millisecond)})},
		{Action: "pop", ExpectedError: ErrNotAvailable},
		{Action: "pop", ExpectedError: ErrNotAvailable, Wait: 50 * time.Millisecond},
		{Action: "pop", ExpectedMessage: "task-3", Wait: 100 * time.Millisecond},
		{Action: "pop", ExpectedError: ErrNotAvailable},
		{Action: "pop", ExpectedMessage: "task-2", Wait: 100 * time.Millisecond},
		{Action: "pop", ExpectedError: ErrNotAvailable},
		{Action: "pop", ExpectedMessage: "task-1", Wait: 100 * time.Millisecond},
		{Action: "pop", ExpectedError: ErrNotFound},
	})
}
