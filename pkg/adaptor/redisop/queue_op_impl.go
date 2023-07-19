package redisop

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/toysmars/distlib-go/pkg/task"
	"github.com/toysmars/distlib-go/pkg/task/queue"
)

func (op *redisOp) Push(ctx context.Context, item *queue.Item) error {
	buf, err := json.Marshal(item)
	if err != nil {
		return errors.Wrap(queue.ErrMarshal, err.Error())
	}

	// Add the item value to the item set.
	err = op.cmdable.Set(ctx, op.getQueueItemKey(item.Task.Group, item.Id), string(buf), 0).Err()
	if err != nil {
		return errors.Wrap(queue.ErrInternal, err.Error())
	}

	// Add the item to the sorted set.
	err = op.cmdable.ZAdd(ctx, op.getQueueKey(item.Task.Group), redis.Z{
		Score:  float64(item.Score),
		Member: item.Id,
	}).Err()
	if err != nil {
		return errors.Wrap(queue.ErrInternal, err.Error())
	}

	return nil
}

func (op *redisOp) Pop(ctx context.Context, group task.Group) (*queue.Item, error) {
	// Pop the item with highest priority from the sorted set.
	res, err := op.cmdable.ZPopMin(ctx, op.getQueueKey(group), 1).Result()
	if err != nil {
		return nil, errors.Wrap(queue.ErrInternal, err.Error())
	}
	if len(res) == 0 {
		return nil, queue.ErrNotFound
	}

	// Fetch and remove.
	return op.fetchAndRemove(ctx, group, res[0].Member.(string), false)
}

func (op *redisOp) PopScheduled(ctx context.Context, group task.Group) (*queue.Item, error) {
	for {
		cnt, err := op.cmdable.ZCount(ctx, op.getQueueKey(group), "-inf", "+inf").Result()
		if err != nil {
			return nil, errors.Wrap(queue.ErrInternal, err.Error())
		}
		if cnt == 0 {
			return nil, queue.ErrNotFound
		}
		// Get the item that passed the scheduled time.
		ids, err := op.cmdable.ZRangeByScore(ctx, op.getQueueKey(group), &redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.FormatInt(time.Now().UnixMilli(), 10),
			Offset: 0,
			Count:  1,
		}).Result()
		if err != nil {
			return nil, errors.Wrap(queue.ErrInternal, err.Error())
		}
		if len(ids) == 0 {
			return nil, queue.ErrNotAvailable
		}

		// Fetch and remove.
		item, err := op.fetchAndRemove(ctx, group, ids[0], true)
		// If it failed to fetch, the item is fetched by other worker.
		if errors.Is(err, queue.ErrNotFound) {
			continue
		}
		return item, err
	}
}

func (op *redisOp) fetchAndRemove(ctx context.Context, group task.Group, id string, mustRemove bool) (*queue.Item, error) {
	// Get the item from the item set.
	buf, err := op.cmdable.Get(ctx, op.getQueueItemKey(group, id)).Result()
	if err != nil {
		return nil, errors.Wrap(queue.ErrInternal, err.Error())
	}

	// Unmarshal the item.
	var item queue.Item
	err = json.Unmarshal([]byte(buf), &item)
	if err != nil {
		return nil, errors.Wrap(queue.ErrMarshal, err.Error())
	}

	// Remove the item from the sorted set.
	removed, err := op.cmdable.ZRem(ctx, op.getQueueKey(group), item.Id).Result()
	if err != nil {
		return nil, errors.Wrap(queue.ErrInternal, err.Error())
	}
	if removed == 0 && mustRemove {
		return nil, queue.ErrNotFound
	}

	// Remove the item from the item set.
	removed, err = op.cmdable.Del(ctx, op.getQueueItemKey(group, item.Id)).Result()
	if err != nil {
		return nil, errors.Wrap(queue.ErrInternal, err.Error())
	}
	if removed == 0 && mustRemove {
		return nil, queue.ErrNotFound
	}

	return &item, nil
}

func (op *redisOp) getQueueKey(group task.Group) string {
	return fmt.Sprintf("%s:%s:task-queue", group.Namespace, group.Name)
}

func (op *redisOp) getQueueItemKey(group task.Group, id string) string {
	return fmt.Sprintf("%s:%s:task-queue:%s", group.Namespace, group.Name, id)
}
