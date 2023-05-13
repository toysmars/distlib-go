package inmemory

import "sync"

func NewInMemoryOperator() Operator {
	return &inMemoryOp{}
}

type inMemoryOp struct {
	pq priorityQueue
	mu sync.Mutex
}
