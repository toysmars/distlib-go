package servicediscovery

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

func NewServicConsumer(option ServiceConsumerOption) (ServiceConsumer, error) {
	if option.DiscoveryDuration <= 0 {
		return nil, errors.Wrap(ErrInit, "invalid DiscoveryDuration")
	}
	if option.ServiceConsumerOperator == nil {
		return nil, errors.Wrap(ErrInit, "invalid ServiceProviderOperator")
	}

	return &serviceConsumerImpl{
		option: option,
		op:     option.ServiceConsumerOperator,
	}, nil
}

type serviceConsumerImpl struct {
	option           ServiceConsumerOption
	op               ServiceConsumerOperator
	muRun            sync.Mutex
	muDiscover       sync.Mutex
	healthStatusList atomic.Value
	discoveryErr     atomic.Value
}

func (s *serviceConsumerImpl) List(ctx context.Context) ([]HealthStatus, error) {
	statusList := s.healthStatusList.Load().([]HealthStatus)
	discoveryErr := s.discoveryErr.Load().(error)

	return statusList, discoveryErr
}

func (s *serviceConsumerImpl) Run(ctx context.Context) error {
	if !s.muRun.TryLock() {
		return ErrDuplicateRun
	}
	defer s.muRun.Unlock()

	go func() {
		ticker := time.NewTicker(s.option.DiscoveryDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.discover(ctx)
			}
		}
	}()
	return nil
}

func (s *serviceConsumerImpl) discover(ctx context.Context) {
	if !s.muDiscover.TryLock() {
		return
	}
	defer s.muDiscover.Unlock()

	// Get the latest health status from the service registry.
	newList, err := s.op.ListHealthStatus(ctx, s.option.Identity)
	if err != nil {
		s.discoveryErr.Store(err)
		if s.option.DiscoveryErrCallback != nil {
			s.option.DiscoveryErrCallback(err)
		}
		return
	}

	// Update available stauts
	now := time.Now()
	for i := 0; i < len(newList); i++ {
		newList[i].IsAvailble = newList[i].AvailableUntil.After(now)
	}

	// Sort the new serice health status list.
	sort.Slice(newList, func(i int, j int) bool {
		return newList[i].Identity.Instance < newList[j].Identity.Instance
	})

	// Get the diff from the previous health check.
	oldList := s.healthStatusList.Load().([]HealthStatus)
	becameAvailableList, becameUnavailableList := getDiff(oldList, newList)

	// Invoke the callback for service health status changes.
	if s.option.DiscoveryCallback != nil {
		for _, st := range becameUnavailableList {
			s.option.DiscoveryCallback(st)
		}
		for _, st := range becameAvailableList {
			s.option.DiscoveryCallback(st)
		}
	}

	// Clean up the registry
	for _, st := range becameUnavailableList {
		s.op.DeregisterHealthStatus(ctx, st.Identity)
	}

	// Store the latest service status.
	s.healthStatusList.Store(newList)
}

func getDiff(oldList, newList []HealthStatus) (becameAvailable, becameUnavailable []HealthStatus) {
	// Invariant - old and new are sorted.
	i := 0
	j := 0
	for i < len(oldList) && j < len(newList) {
		old := &oldList[i]
		new := &newList[j]
		if old.Identity.Instance < new.Identity.Instance {
			if old.IsAvailble {
				becameUnavailable = append(becameUnavailable, *old)
			}
			i++
		} else if old.Identity.Instance > new.Identity.Instance {
			if new.IsAvailble {
				becameAvailable = append(becameAvailable, *new)
			}
			j++
		} else {
			if !old.IsAvailble && new.IsAvailble {
				becameAvailable = append(becameAvailable, *new)
			} else if old.IsAvailble && !new.IsAvailble {
				becameUnavailable = append(becameUnavailable, *old)
			}
			i++
			j++
		}
	}
	for i < len(oldList) {
		if oldList[i].IsAvailble {
			becameUnavailable = append(becameUnavailable, oldList[i])
		}
		i++
	}
	for j < len(newList) {
		if newList[j].IsAvailble {
			becameAvailable = append(becameAvailable, newList[j])
		}
		j++
	}
	return
}
