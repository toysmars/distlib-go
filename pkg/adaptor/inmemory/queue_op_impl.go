package inmemory

import (
	"container/heap"
	"context"
	"time"

	"github.com/toysmars/distlib-go/pkg/task"
	"github.com/toysmars/distlib-go/pkg/task/queue"
)

func (op *inMemoryOp) Push(ctx context.Context, item *queue.Item) error {
	op.mu.Lock()
	defer op.mu.Unlock()

	heap.Push(&op.pq, item)
	return nil
}

func (op *inMemoryOp) Pop(ctx context.Context, group task.Group) (*queue.Item, error) {
	op.mu.Lock()
	defer op.mu.Unlock()

	if len(op.pq) == 0 {
		return nil, queue.ErrNotFound
	}

	item := op.pq[0]
	heap.Pop(&op.pq)
	return item, nil
}

func (op *inMemoryOp) PopScheduled(ctx context.Context, group task.Group) (*queue.Item, error) {
	op.mu.Lock()
	defer op.mu.Unlock()

	if len(op.pq) == 0 {
		return nil, queue.ErrNotFound
	}

	item := op.pq[0]
	// See if the item is available per schedule.
	scheduledFor := time.UnixMilli(item.Score)
	if time.Now().Before(scheduledFor) {
		return nil, queue.ErrNotAvailable
	}
	heap.Pop(&op.pq)
	return item, nil
}

type priorityQueue []*queue.Item

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
	item := x.(*queue.Item)
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
