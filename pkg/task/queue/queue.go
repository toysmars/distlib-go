package queue

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/toysmars/distlib-go/pkg/task"
)

var (
	ErrInit         = errors.New("init error")
	ErrQueueFull    = errors.New("queue is full")
	ErrNotFound     = errors.New("not found")
	ErrNotAvailable = errors.New("item not available")
	ErrMarshal      = errors.New("marshal error")
	ErrInternal     = errors.New("internal error")
)

// Kind is the kind of the queue.
type Kind int

const (
	Fifo Kind = iota
	Priority
	Scheduled
)

// Option specifies the option for a queue.
type Option struct {
	// TaskGroup identifies the task group.
	TaskGroup task.Group
	// Kind is the kind of the queue.
	Kind Kind
	// PollInterval is the desired wait period between consecutive `Pop()` function calls for available item polling.
	PollInterval time.Duration
	// Factory is the Operator factory.
	Factory OperatorFactory
}

// Item represents an item in the queue.
type Item struct {
	Id   string
	Task *task.Task
	// Score is the priority score of the item. The lower the socre, the higher the priority.
	// It is desired to fetch items with the lowest score first.
	// The score is determined per queue kind as following:
	//  - Fifo: Incremental sequence number.
	//  - Priority: The priority of the task.
	//  - Scheduled: Epoch in millisecond of the scheduled time.
	Score    int64
	QueuedAt time.Time
}

// Queue is the task queue.
type Queue interface {
	// Push pushes a task into the queue.
	Push(ctx context.Context, task *task.Task) (*Item, error)

	// Pop pops the Item object with the highest priority.
	//  it returns ErrNoItem when there's no items in the queue.
	//  it returns ErrNotAvailable when there's no available items in the queue(the scheduled time not yet reached).
	Pop(ctx context.Context) (*Item, error)

	// TryPop pops the Item object with the highest priority.
	// This function can wait up to `timeout` until there's an item become available.
	//  it returns ErrNoItem when there's no items in the queue.
	//  it returns ErrNotAvailable when there's no available items in the queue(the scheduled time not yet reached).
	TryPop(ctx context.Context, timeout time.Duration) (*Item, error)
}

// OperatorFactory is the factory for Operator.
type OperatorFactory interface {
	CreateOperator(option Option) Operator
}

// Operator is the interface should be implmented by the underlying queue implementation.
type Operator interface {
	// Push pushes a task item into the queue.
	Push(ctx context.Context, item *Item) error

	// Pop pops the Item object with the highest priority.
	//  it returns ErrNoItem when there's no items in the queue.
	//  it returns ErrNotAvailable when there's no available items in the queue(the scheduled time not yet reached).
	Pop(ctx context.Context) (*Item, error)
}
