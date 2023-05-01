package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisOperatorFactory(client *redis.Client) OperatorFactory {
	return &redisOpFactory{
		client: client,
	}
}

type redisOpFactory struct {
	client *redis.Client
}

func (f *redisOpFactory) CreateOperator(option Option) Operator {
	return NewRedisOperator(f.client, option)
}

func NewRedisOperator(client *redis.Client, option Option) Operator {
	return &redisOp{
		client: client,
		option: option,
	}
}

type redisOp struct {
	client *redis.Client
	option Option
}

func (op *redisOp) Push(ctx context.Context, item *Item) error {
	buf, err := json.Marshal(item)
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}

	// Add the item value to the item set.
	err = op.client.Set(ctx, op.getQueueItemKey(item.Id), string(buf), 0).Err()
	if err != nil {
		return errors.Wrap(ErrInternal, err.Error())
	}

	// Add the item to the sorted set.
	err = op.client.ZAdd(ctx, op.getQueueKey(), redis.Z{
		Score:  float64(item.Score),
		Member: item.Id,
	}).Err()
	if err != nil {
		return errors.Wrap(ErrInternal, err.Error())
	}

	return nil
}

func (op *redisOp) Pop(ctx context.Context) (*Item, error) {
	if op.option.Kind == Scheduled {
		return op.popScheduled(ctx)
	}
	return op.popHead(ctx)
}

func (op *redisOp) popHead(ctx context.Context) (*Item, error) {
	// Pop the item with highest priority from the sorted set.
	res, err := op.client.ZPopMin(ctx, op.getQueueKey(), 1).Result()
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}

	// Fetch and remove.
	return op.fetchAndRemove(ctx, res[0].Member.(string), false)
}

func (op *redisOp) popScheduled(ctx context.Context) (*Item, error) {
	for {
		cnt, err := op.client.ZCount(ctx, op.getQueueKey(), "-inf", "+inf").Result()
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
		if cnt == 0 {
			return nil, ErrNotFound
		}
		// Get the item that passed the scheduled time.
		ids, err := op.client.ZRangeByScore(ctx, op.getQueueKey(), &redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.FormatInt(time.Now().UnixMilli(), 10),
			Offset: 0,
			Count:  1,
		}).Result()
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
		if len(ids) == 0 {
			return nil, ErrNotAvailable
		}

		// Fetch and remove.
		item, err := op.fetchAndRemove(ctx, ids[0], true)
		// If it failed to fetch, the item is fetched by other worker.
		if errors.Is(err, ErrNotFound) {
			continue
		}
		return item, err
	}
}

func (op *redisOp) fetchAndRemove(ctx context.Context, id string, mustRemove bool) (*Item, error) {
	// Get the item from the item set.
	buf, err := op.client.Get(ctx, op.getQueueItemKey(id)).Result()
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	// Unmarshal the item.
	var item Item
	err = json.Unmarshal([]byte(buf), &item)
	if err != nil {
		return nil, errors.Wrap(ErrMarshal, err.Error())
	}

	// Remove the item from the sorted set.
	removed, err := op.client.ZRem(ctx, op.getQueueKey(), item.Id).Result()
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	if removed == 0 && mustRemove {
		return nil, ErrNotFound
	}

	// Remove the item from the item set.
	removed, err = op.client.Del(ctx, op.getQueueItemKey(item.Id)).Result()
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	if removed == 0 && mustRemove {
		return nil, ErrNotFound
	}

	return &item, nil
}

func (op *redisOp) getQueueKey() string {
	return fmt.Sprintf("%s:%s:task-queue", op.option.TaskGroup.Namespace, op.option.TaskGroup.Name)
}

func (op *redisOp) getQueueItemKey(id string) string {
	return fmt.Sprintf("%s:%s:task-queue:%s", op.option.TaskGroup.Namespace, op.option.TaskGroup.Name, id)
}
