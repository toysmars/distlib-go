package queue

import (
	context "context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toysmars/distlib-go/pkg/task"
)

type queueTestSetup struct {
	ctx           context.Context
	option        Option
	queue         *queueImpl
	mockOp        *MockOperator
	mockOpFactory *MockOperatorFactory
}

func createQueueTestSetup(t *testing.T, kind Kind) queueTestSetup {
	t.Helper()

	factory := NewMockOperatorFactory(t)
	op := NewMockOperator(t)
	option := Option{
		Kind:         kind,
		Factory:      factory,
		PollInterval: 10 * time.Millisecond,
	}

	factory.On("CreateOperator", option).Return(op)
	queue, err := New(option)
	assert.NoError(t, err)

	return queueTestSetup{
		ctx:           context.Background(),
		option:        option,
		mockOp:        op,
		mockOpFactory: factory,
		queue:         queue.(*queueImpl),
	}
}

func TestQueue_Push(t *testing.T) {
	t.Parallel()

	testTask := &task.Task{
		Message: "my test task",
	}

	t.Run("test queue.Push() success", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		ts.mockOp.On("Push", ts.ctx, mock.Anything).Return(nil)

		item, err := ts.queue.Push(ts.ctx, testTask)

		assert.NoError(t, err)
		assert.Equal(t, testTask, item.Task)
	})
	t.Run("test queue.Push() fail with ErrInternal", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		ts.mockOp.On("Push", ts.ctx, mock.Anything).Return(ErrInternal)

		item, err := ts.queue.Push(ts.ctx, testTask)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInternal))
		assert.Nil(t, item)
	})
}

func TestQueue_Pop(t *testing.T) {
	t.Parallel()

	testItem := &Item{
		Task: &task.Task{
			Message: "my test task",
		},
	}

	testWithExpectedError := func(t *testing.T, ts *queueTestSetup, expectedErr error) {
		ts.mockOp.On("Pop", ts.ctx).Return(nil, expectedErr)

		item, err := ts.queue.Pop(ts.ctx)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, item)
	}

	t.Run("test queue.Pop() success", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		ts.mockOp.On("Pop", ts.ctx).Return(testItem, nil)

		item, err := ts.queue.Pop(ts.ctx)

		assert.NoError(t, err)
		assert.Equal(t, testItem.Task, item.Task)
	})
	t.Run("test queue.Pop() fail with ErrNotFound", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrNotFound)
	})
	t.Run("test queue.Pop() fail with ErrNotAvailable", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrNotAvailable)
	})
	t.Run("test queue.Pop() fail with ErrInternal", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrInternal)
	})
}

func TestQueue_TryPop(t *testing.T) {
	t.Parallel()

	testItem := &Item{
		Task: &task.Task{
			Message: "my test task",
		},
	}

	testWithExpectedError := func(t *testing.T, ts *queueTestSetup, expectedErr error) {
		ts.mockOp.On("Pop", ts.ctx).Return(nil, expectedErr)

		item, err := ts.queue.TryPop(ts.ctx, 100*time.Millisecond)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, item)
	}
	getNumberOfPopCalles := func(ts *queueTestSetup) int {
		calles := 0
		for _, call := range ts.mockOp.Calls {
			if call.Method == "Pop" {
				calles++
			}
		}
		return calles
	}

	t.Run("test queue.TryPop() success", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		ts.mockOp.On("Pop", ts.ctx).Return(testItem, nil)

		item, err := ts.queue.TryPop(ts.ctx, 100*time.Millisecond)

		assert.NoError(t, err)
		assert.Equal(t, testItem.Task, item.Task)
		ts.mockOp.AssertNumberOfCalls(t, "Pop", 1)
	})
	t.Run("test queue.TryPop() success with wating", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		called := 0
		ts.mockOp.On("Pop", ts.ctx).Return(
			func(ctx context.Context) *Item {
				called++
				if called == 5 {
					return testItem
				}
				return nil
			},
			func(ctx context.Context) error {
				if called < 5 {
					return ErrNotAvailable
				}
				if called == 5 {
					return nil
				}
				return ErrNotFound
			},
		)

		item, err := ts.queue.TryPop(ts.ctx, 100*time.Millisecond)

		assert.NoError(t, err)
		assert.Equal(t, testItem.Task, item.Task)
		ts.mockOp.AssertNumberOfCalls(t, "Pop", 5)
	})
	t.Run("test queue.Pop() fail with ErrNotFound", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrNotFound)

		assert.True(t, getNumberOfPopCalles(&ts) > 5)
	})
	t.Run("test queue.Pop() fail with ErrNotAvailable", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrNotAvailable)

		assert.True(t, getNumberOfPopCalles(&ts) > 5)
	})
	t.Run("test queue.Pop() fail with ErrInternal", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)
		testWithExpectedError(t, &ts, ErrInternal)

		ts.mockOp.AssertNumberOfCalls(t, "Pop", 1)
	})
}

func TestQueue_createItem(t *testing.T) {
	testTask := &task.Task{
		Message: "my test task",
	}
	t.Run("test queueImpl.createItem() with Fifo", func(t *testing.T) {
		ts := createQueueTestSetup(t, Fifo)

		item1 := ts.queue.createItem(testTask)
		item2 := ts.queue.createItem(testTask)
		item3 := ts.queue.createItem(testTask)
		item4 := ts.queue.createItem(testTask)
		item5 := ts.queue.createItem(testTask)

		assert.True(t, item1.Score < item2.Score)
		assert.True(t, item2.Score < item3.Score)
		assert.True(t, item3.Score < item4.Score)
		assert.True(t, item4.Score < item5.Score)
	})
	t.Run("test queueImpl.createItem() with Priority", func(t *testing.T) {
		ts := createQueueTestSetup(t, Priority)

		item1 := ts.queue.createItem(&task.Task{Option: task.Option{Priority: 5}})
		item2 := ts.queue.createItem(&task.Task{Option: task.Option{Priority: 4}})
		item3 := ts.queue.createItem(&task.Task{Option: task.Option{Priority: 3}})
		item4 := ts.queue.createItem(&task.Task{Option: task.Option{Priority: 2}})
		item5 := ts.queue.createItem(&task.Task{Option: task.Option{Priority: 10}})

		assert.Equal(t, int64(5), item1.Score)
		assert.Equal(t, int64(4), item2.Score)
		assert.Equal(t, int64(3), item3.Score)
		assert.Equal(t, int64(2), item4.Score)
		assert.Equal(t, int64(10), item5.Score)
	})
	t.Run("test queueImpl.createItem() with Scheduled", func(t *testing.T) {
		ts := createQueueTestSetup(t, Scheduled)

		time1 := time.Now().Add(5 * time.Minute)
		item1 := ts.queue.createItem(&task.Task{Option: task.Option{ScheduledFor: time1}})
		time2 := time.Now().Add(15 * time.Minute)
		item2 := ts.queue.createItem(&task.Task{Option: task.Option{ScheduledFor: time2}})
		time3 := time.Now().Add(500 * time.Millisecond)
		item3 := ts.queue.createItem(&task.Task{Option: task.Option{ScheduledFor: time3}})
		time4 := time.Now().Add(7 * time.Hour)
		item4 := ts.queue.createItem(&task.Task{Option: task.Option{ScheduledFor: time4}})
		time5 := time.Time{}
		item5 := ts.queue.createItem(&task.Task{Option: task.Option{ScheduledFor: time5}})

		assert.Equal(t, time1.UnixMilli(), item1.Score)
		assert.Equal(t, time2.UnixMilli(), item2.Score)
		assert.Equal(t, time3.UnixMilli(), item3.Score)
		assert.Equal(t, time4.UnixMilli(), item4.Score)
		assert.Equal(t, time5.UnixMilli(), item5.Score)
	})
}
