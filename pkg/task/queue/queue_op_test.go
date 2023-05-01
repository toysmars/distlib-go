package queue

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/toysmars/distlib-go/pkg/task"
)

type operatorTestSetup struct {
	ctx    context.Context
	op     Operator
	option Option
}

type testAction struct {
	Action          string
	Task            *task.Task
	Wait            time.Duration
	ExpectedError   error
	ExpectedMessage string
}

func createOperatorTestSetup(op Operator, option Option) operatorTestSetup {
	return operatorTestSetup{
		ctx:    context.Background(),
		op:     op,
		option: option,
	}
}

func (ts *operatorTestSetup) createTask(message string, option task.Option) *task.Task {
	return &task.Task{
		Message: message,
		Group:   ts.option.TaskGroup,
		Option:  option,
	}
}

func (ts *operatorTestSetup) createItem(task *task.Task, score int64) *Item {
	return &Item{
		Id:       uuid.NewString(),
		Task:     task,
		Score:    score,
		QueuedAt: time.Now(),
	}
}

func (ts *operatorTestSetup) runTestActions(t *testing.T, actions []testAction) {
	t.Helper()

	for i, action := range actions {
		var item *Item
		var err error
		getScore := func() int64 {
			switch ts.option.Kind {
			case Priority:
				return int64(action.Task.Option.Priority)
			case Scheduled:
				return action.Task.Option.ScheduledFor.UnixMilli()
			}
			return int64(i)
		}
		if action.Wait > 0 {
			time.Sleep(action.Wait)
		}
		switch action.Action {
		case "push":
			err = ts.op.Push(ts.ctx, ts.createItem(action.Task, getScore()))
		case "pop":
			item, err = ts.op.Pop(ts.ctx)
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
