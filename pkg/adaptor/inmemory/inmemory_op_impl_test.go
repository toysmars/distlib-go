package inmemory

import (
	"testing"

	queueTest "github.com/toysmars/distlib-go/pkg/task/queue/test"
)

func TestInMemoryQueueOperator(t *testing.T) {
	t.Parallel()

	queueTest.RunOperatorUnitTests(t, NewInMemoryOperator())
}
