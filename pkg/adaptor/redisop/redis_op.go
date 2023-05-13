package redisop

import (
	"github.com/toysmars/distlib-go/pkg/servicediscovery"
	"github.com/toysmars/distlib-go/pkg/task/queue"
)

// Operator implements:
//   - queue.Operator
//   - servicediscovery.Operator
type Operator interface {
	queue.Operator
	servicediscovery.Operator
}
