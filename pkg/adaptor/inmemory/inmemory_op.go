package inmemory

import "github.com/toysmars/distlib-go/pkg/task/queue"

// Operator implements:
//   - queue.Operator
type Operator interface {
	queue.Operator
}
