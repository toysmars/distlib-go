package servicediscovery

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInit         = errors.New("init error")
	ErrDuplicateRun = errors.New("duplicate runs")
	ErrInternal     = errors.New("internal error")
	ErrMarshal      = errors.New("marshal error")
)

// Identity is the service identity managed by the service discovery.
type Identity struct {
	Namespace string
	Service   string
	Instance  string
}

// HealthStatus is the service status.
type HealthStatus struct {
	Identity       Identity
	IsAvailble     bool
	LastHeartbeat  time.Time
	AvailableUntil time.Time
}

// ServiceProviderOption defines the behavior of the service provider.
type ServiceProviderOption struct {
	Identity                Identity
	HeartbeatDuration       time.Duration
	AvailableDuration       time.Duration
	RegisterErrCallback     RegisterErrCallback
	ServiceProviderOperator ServiceProviderOperator
}

type RegisterErrCallback func(err error)

// ServiceProvider provides a service via the service discovery.
type ServiceProvider interface {
	Run(ctx context.Context) error
}

// ServiceConsumerOption defines the behavior of the service consumer.
type ServiceConsumerOption struct {
	Identity                Identity
	DiscoveryDuration       time.Duration
	DiscoveryCallback       DiscoveryCallback
	DiscoveryErrCallback    DiscoveryErrCallback
	ServiceConsumerOperator ServiceConsumerOperator
}

type DiscoveryCallback func(status HealthStatus)
type DiscoveryErrCallback func(err error)

// ServiceConsumer consumes services via the service discovery.
type ServiceConsumer interface {
	List(ctx context.Context) ([]HealthStatus, error)
	Run(ctx context.Context) error
}

// ServiceProviderOperator is the interface that should be implmented by the underlying implementation for ServiceProvider.
type ServiceProviderOperator interface {
	RegisterHealthStatus(ctx context.Context, status HealthStatus) error
}

// ServiceConsumerOperator is the interface that should be implmented by the underlying implementation for ServiceConsume.
type ServiceConsumerOperator interface {
	ListHealthStatus(ctx context.Context, identity Identity) ([]HealthStatus, error)
	DeregisterHealthStatus(ctx context.Context, identity Identity) error
}

// Operator is the interface that should be implmented by the underlying implementation for both ServiceProvider and ServiceConsume.
type Operator interface {
	ServiceProviderOperator
	ServiceConsumerOperator
}
