package queue

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

func NewInMemoryOperatorFactory() OperatorFactory {
	return &inMemoryOpFactory{}
}

type inMemoryOpFactory struct{}

func (f *inMemoryOpFactory) CreateOperator(option Option) Operator {
	return NewInMemoryOperator(option)
}

func NewInMemoryOperator(option Option) Operator {
	return &inMemoryOp{
		option: option,
	}
}

type inMemoryOp struct {
	option Option
	pq     priorityQueue
	mu     sync.Mutex
}

func (op *inMemoryOp) Push(ctx context.Context, item *Item) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	heap.Push(&op.pq, item)
	return nil
}

func (op *inMemoryOp) Pop(ctx context.Context) (*Item, error) {
	op.mu.Lock()
	defer op.mu.Unlock()

	if len(op.pq) == 0 {
		return nil, ErrNotFound
	}

	item := op.pq[0]
	// See if the item is available per schedule.
	if op.option.Kind == Scheduled && time.Now().Before(item.Task.Option.ScheduledFor) {
		return nil, ErrNotAvailable
	}
	heap.Pop(&op.pq)
	return item, nil
}

type priorityQueue []*Item

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Score < pq[j].Score
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
