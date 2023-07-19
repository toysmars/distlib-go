package redisop

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	sd "github.com/toysmars/distlib-go/pkg/servicediscovery"
)

func (op *redisOp) RegisterHealthStatus(ctx context.Context, status sd.HealthStatus) error {
	// Registry key.
	key := getHealthStatusRegistryKey(status.Identity)

	// Marshall the health status.
	buf, err := json.Marshal(&status)
	if err != nil {
		return errors.Wrap(sd.ErrMarshal, err.Error())
	}

	// Add the health status to the registry map.
	err = op.cmdable.HSet(ctx, key, status.Identity.Instance, string(buf)).Err()
	if err != nil {
		return errors.Wrap(sd.ErrInternal, err.Error())
	}

	return nil
}

func (op *redisOp) ListHealthStatus(ctx context.Context, identity sd.Identity) ([]sd.HealthStatus, error) {
	// Registry key.
	key := getHealthStatusRegistryKey(identity)

	// Get the members in the registry map.
	statusMap, err := op.cmdable.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(sd.ErrInternal, err.Error())
	}

	// Unmarshal the health status.
	var result []sd.HealthStatus
	for _, statusBuf := range statusMap {
		var status sd.HealthStatus
		err = json.Unmarshal([]byte(statusBuf), &status)
		if err != nil {
			return nil, errors.Wrap(sd.ErrMarshal, err.Error())
		}
		result = append(result, status)
	}

	return result, nil
}

func (op *redisOp) DeregisterHealthStatus(ctx context.Context, identity sd.Identity) error {
	// Registry key.
	key := getHealthStatusRegistryKey(identity)

	// Delete the health status from the registry map.
	err := op.cmdable.HDel(ctx, key, identity.Instance).Err()
	if err != nil {
		return errors.Wrap(sd.ErrInternal, err.Error())
	}
	return nil
}

func getHealthStatusRegistryKey(identity sd.Identity) string {
	return fmt.Sprintf("%s:%s:health-status", identity.Namespace, identity.Service)
}
