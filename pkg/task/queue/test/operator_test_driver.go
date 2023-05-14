package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/toysmars/distlib-go/pkg/task"
	"github.com/toysmars/distlib-go/pkg/task/queue"
)

type operatorTestDriver struct {
	ctx    context.Context
	op     queue.Operator
	option queue.Option
}

type operatorTestAction struct {
	Action          string
	Task            *task.Task
	Wait            time.Duration
	ExpectedError   error
	ExpectedMessage string
}

func createOperatorTestDriver(op queue.Operator, option queue.Option) operatorTestDriver {
	return operatorTestDriver{
		ctx:    context.Background(),
		op:     op,
		option: option,
	}
}

func (ts *operatorTestDriver) CreateTask(message string, option task.Option) *task.Task {
	return &task.Task{
		Message: message,
		Group:   ts.option.TaskGroup,
		Option:  option,
	}
}

func (ts *operatorTestDriver) CreateItem(task *task.Task, score int64) *queue.Item {
	return &queue.Item{
		Id:       uuid.NewString(),
		Task:     task,
		Score:    score,
		QueuedAt: time.Now(),
	}
}

func (ts *operatorTestDriver) RunTestActions(t *testing.T, actions []operatorTestAction) {
	t.Helper()

	for i, action := range actions {
		var item *queue.Item
		var err error
		getScore := func() int64 {
			switch ts.option.Kind {
			case queue.Priority:
				return int64(action.Task.Option.Priority)
			case queue.Scheduled:
				return action.Task.Option.ScheduledFor.UnixMilli()
			}
			return int64(i)
		}
		if action.Wait > 0 {
			time.Sleep(action.Wait)
		}
		switch action.Action {
		case "push":
			err = ts.op.Push(ts.ctx, ts.CreateItem(action.Task, getScore()))
		case "pop":
			if ts.option.Kind == queue.Scheduled {
				item, err = ts.op.PopScheduled(ts.ctx, ts.option.TaskGroup)
			} else {
				item, err = ts.op.Pop(ts.ctx, ts.option.TaskGroup)
			}
		}
		if action.ExpectedError == nil {
			assert.NoError(t, err, action)
		} else {
			assert.True(t, errors.Is(err, action.ExpectedError), action)
		}
		if action.ExpectedMessage != "" {
			assert.Equal(t, action.ExpectedMessage, item.Task.Message, action)
		}
	}
}

func RunOperatorUnitTests(t *testing.T, op queue.Operator) {
	t.Run("queue operator test - fifo", func(t *testing.T) {
		option := queue.Option{
			Kind: queue.Fifo,
			TaskGroup: task.Group{
				Namespace: "unit-test",
				Name:      "test-queue-operator-fifo",
			},
		}
		ts := createOperatorTestDriver(op, option)
		ts.RunTestActions(t, []operatorTestAction{
			{Action: "pop", ExpectedError: queue.ErrNotFound},
			{Action: "push", Task: ts.CreateTask("task-1", task.Option{})},
			{Action: "push", Task: ts.CreateTask("task-2", task.Option{})},
			{Action: "pop", ExpectedMessage: "task-1"},
			{Action: "push", Task: ts.CreateTask("task-3", task.Option{})},
			{Action: "pop", ExpectedMessage: "task-2"},
			{Action: "pop", ExpectedMessage: "task-3"},
			{Action: "pop", ExpectedError: queue.ErrNotFound},
			{Action: "push", Task: ts.CreateTask("task-4", task.Option{})},
			{Action: "push", Task: ts.CreateTask("task-5", task.Option{})},
			{Action: "pop", ExpectedMessage: "task-4"},
			{Action: "push", Task: ts.CreateTask("task-6", task.Option{})},
			{Action: "pop", ExpectedMessage: "task-5"},
			{Action: "pop", ExpectedMessage: "task-6"},
			{Action: "pop", ExpectedError: queue.ErrNotFound},
		})
	})
	t.Run("queue operator test - priority", func(t *testing.T) {
		option := queue.Option{
			Kind: queue.Priority,
			TaskGroup: task.Group{
				Namespace: "unit-test",
				Name:      "test-queue-operator-priority",
			},
		}
		ts := createOperatorTestDriver(op, option)
		ts.RunTestActions(t, []operatorTestAction{
			{Action: "pop", ExpectedError: queue.ErrNotFound},
			{Action: "push", Task: ts.CreateTask("task-1", task.Option{Priority: 3})},
			{Action: "push", Task: ts.CreateTask("task-2", task.Option{Priority: 5})},
			{Action: "push", Task: ts.CreateTask("task-3", task.Option{Priority: 2})},
			{Action: "pop", ExpectedMessage: "task-3"},
			{Action: "push", Task: ts.CreateTask("task-4", task.Option{Priority: 4})},
			{Action: "pop", ExpectedMessage: "task-1"},
			{Action: "push", Task: ts.CreateTask("task-5", task.Option{Priority: 1})},
			{Action: "pop", ExpectedMessage: "task-5"},
			{Action: "pop", ExpectedMessage: "task-4"},
			{Action: "pop", ExpectedMessage: "task-2"},
			{Action: "pop", ExpectedError: queue.ErrNotFound},
		})
	})
	t.Run("queue operator test - scheduled", func(t *testing.T) {
		now := time.Now()
		option := queue.Option{
			Kind: queue.Scheduled,
			TaskGroup: task.Group{
				Namespace: "unit-test",
				Name:      "test-queue-operator-scheduled",
			},
		}
		ts := createOperatorTestDriver(op, option)
		ts.RunTestActions(t, []operatorTestAction{
			{Action: "pop", ExpectedError: queue.ErrNotFound},
			{Action: "push", Task: ts.CreateTask("task-1", task.Option{ScheduledFor: now.Add(300 * time.Millisecond)})},
			{Action: "push", Task: ts.CreateTask("task-2", task.Option{ScheduledFor: now.Add(200 * time.Millisecond)})},
			{Action: "push", Task: ts.CreateTask("task-3", task.Option{ScheduledFor: now.Add(100 * time.Millisecond)})},
			{Action: "pop", ExpectedError: queue.ErrNotAvailable},
			{Action: "pop", ExpectedError: queue.ErrNotAvailable, Wait: 50 * time.Millisecond},
			{Action: "pop", ExpectedMessage: "task-3", Wait: 100 * time.Millisecond},
			{Action: "pop", ExpectedError: queue.ErrNotAvailable},
			{Action: "pop", ExpectedMessage: "task-2", Wait: 100 * time.Millisecond},
			{Action: "pop", ExpectedError: queue.ErrNotAvailable},
			{Action: "pop", ExpectedMessage: "task-1", Wait: 100 * time.Millisecond},
			{Action: "pop", ExpectedError: queue.ErrNotFound},
		})
	})
}
