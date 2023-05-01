package queue

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/toysmars/distlib-go/pkg/task"
)

func New(option Option) (Queue, error) {
	if option.Kind < Fifo || option.Kind > Scheduled {
		return nil, errors.Wrap(ErrInit, "invalid Kind")
	}
	if option.PollInterval <= 0 {
		return nil, errors.Wrap(ErrInit, "invalid PollInternal")
	}
	if option.Factory == nil {
		return nil, errors.Wrap(ErrInit, "Factory must be supplied")
	}
	return &queueImpl{
		option: option,
		op:     option.Factory.CreateOperator(option),
	}, nil
}

type queueImpl struct {
	option Option
	seq    int64
	op     Operator
}

func (q *queueImpl) Push(ctx context.Context, task *task.Task) (*Item, error) {
	item := q.createItem(task)
	err := q.op.Push(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (q *queueImpl) Pop(ctx context.Context) (*Item, error) {
	return q.op.Pop(ctx)
}

func (q *queueImpl) TryPop(ctx context.Context, timeout time.Duration) (item *Item, err error) {
	startsAt := time.Now()
	for time.Since(startsAt) < timeout {
		item, err = q.op.Pop(ctx)
		if err == nil {
			return item, nil
		}
		if !errors.Is(err, ErrNotFound) && !errors.Is(err, ErrNotAvailable) {
			return nil, err
		}
		time.Sleep(q.option.PollInterval)
	}
	return nil, err
}

func (q *queueImpl) createItem(task *task.Task) *Item {
	const maxSeq = int64(1)<<54 - 1
	getScore := func() int64 {
		switch q.option.Kind {
		case Priority:
			return int64(task.Option.Priority)
		case Scheduled:
			return task.Option.ScheduledFor.UnixMilli()
		default:
			q.seq++
			if q.seq > maxSeq {
				q.seq = 0
			}
			return q.seq
		}
	}
	return &Item{
		Id:       uuid.NewString(),
		Task:     task,
		Score:    getScore(),
		QueuedAt: time.Now(),
	}
}
