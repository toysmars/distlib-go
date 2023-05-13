package servicediscovery

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

func NewServiceProvider(option ServiceProviderOption) (ServiceProvider, error) {
	if option.HeartbeatDuration <= 0 {
		return nil, errors.Wrap(ErrInit, "invalid HeartbeatDuration")
	}
	if option.AvailableDuration < option.HeartbeatDuration {
		return nil, errors.Wrap(ErrInit, "invalid AvailableDuration")
	}
	if option.ServiceProviderOperator == nil {
		return nil, errors.Wrap(ErrInit, "invalid ServiceProviderOperator")
	}

	return &serviceProviderImpl{
		option:              option,
		op:                  option.ServiceProviderOperator,
		registerErrCallback: option.RegisterErrCallback,
	}, nil
}

type serviceProviderImpl struct {
	option              ServiceProviderOption
	op                  ServiceProviderOperator
	registerErrCallback RegisterErrCallback
	muRun               sync.Mutex
	muHeartbeat         sync.Mutex
}

func (s *serviceProviderImpl) Run(ctx context.Context) error {
	if !s.muRun.TryLock() {
		return ErrDuplicateRun
	}
	defer s.muRun.Unlock()

	go func() {
		ticker := time.NewTicker(s.option.HeartbeatDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.heartbeat(ctx)
			}
		}
	}()
	return nil
}

func (s *serviceProviderImpl) heartbeat(ctx context.Context) {
	if !s.muHeartbeat.TryLock() {
		return
	}
	defer s.muHeartbeat.Unlock()

	now := time.Now()
	status := HealthStatus{
		Identity:       s.option.Identity,
		IsAvailble:     true,
		LastHeartbeat:  now,
		AvailableUntil: now.Add(s.option.AvailableDuration),
	}

	if err := s.op.RegisterHealthStatus(ctx, status); err != nil && s.registerErrCallback != nil {
		s.registerErrCallback(err)
	}
}
