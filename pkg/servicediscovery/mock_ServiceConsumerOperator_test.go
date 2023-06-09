// Code generated by mockery v2.26.1. DO NOT EDIT.

package servicediscovery

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockServiceConsumerOperator is an autogenerated mock type for the ServiceConsumerOperator type
type MockServiceConsumerOperator struct {
	mock.Mock
}

type MockServiceConsumerOperator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockServiceConsumerOperator) EXPECT() *MockServiceConsumerOperator_Expecter {
	return &MockServiceConsumerOperator_Expecter{mock: &_m.Mock}
}

// DeregisterHealthStatus provides a mock function with given fields: ctx, identity
func (_m *MockServiceConsumerOperator) DeregisterHealthStatus(ctx context.Context, identity Identity) error {
	ret := _m.Called(ctx, identity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Identity) error); ok {
		r0 = rf(ctx, identity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServiceConsumerOperator_DeregisterHealthStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeregisterHealthStatus'
type MockServiceConsumerOperator_DeregisterHealthStatus_Call struct {
	*mock.Call
}

// DeregisterHealthStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - identity Identity
func (_e *MockServiceConsumerOperator_Expecter) DeregisterHealthStatus(ctx interface{}, identity interface{}) *MockServiceConsumerOperator_DeregisterHealthStatus_Call {
	return &MockServiceConsumerOperator_DeregisterHealthStatus_Call{Call: _e.mock.On("DeregisterHealthStatus", ctx, identity)}
}

func (_c *MockServiceConsumerOperator_DeregisterHealthStatus_Call) Run(run func(ctx context.Context, identity Identity)) *MockServiceConsumerOperator_DeregisterHealthStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Identity))
	})
	return _c
}

func (_c *MockServiceConsumerOperator_DeregisterHealthStatus_Call) Return(_a0 error) *MockServiceConsumerOperator_DeregisterHealthStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceConsumerOperator_DeregisterHealthStatus_Call) RunAndReturn(run func(context.Context, Identity) error) *MockServiceConsumerOperator_DeregisterHealthStatus_Call {
	_c.Call.Return(run)
	return _c
}

// ListHealthStatus provides a mock function with given fields: ctx, identity
func (_m *MockServiceConsumerOperator) ListHealthStatus(ctx context.Context, identity Identity) ([]HealthStatus, error) {
	ret := _m.Called(ctx, identity)

	var r0 []HealthStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, Identity) ([]HealthStatus, error)); ok {
		return rf(ctx, identity)
	}
	if rf, ok := ret.Get(0).(func(context.Context, Identity) []HealthStatus); ok {
		r0 = rf(ctx, identity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]HealthStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, Identity) error); ok {
		r1 = rf(ctx, identity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceConsumerOperator_ListHealthStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListHealthStatus'
type MockServiceConsumerOperator_ListHealthStatus_Call struct {
	*mock.Call
}

// ListHealthStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - identity Identity
func (_e *MockServiceConsumerOperator_Expecter) ListHealthStatus(ctx interface{}, identity interface{}) *MockServiceConsumerOperator_ListHealthStatus_Call {
	return &MockServiceConsumerOperator_ListHealthStatus_Call{Call: _e.mock.On("ListHealthStatus", ctx, identity)}
}

func (_c *MockServiceConsumerOperator_ListHealthStatus_Call) Run(run func(ctx context.Context, identity Identity)) *MockServiceConsumerOperator_ListHealthStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Identity))
	})
	return _c
}

func (_c *MockServiceConsumerOperator_ListHealthStatus_Call) Return(_a0 []HealthStatus, _a1 error) *MockServiceConsumerOperator_ListHealthStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceConsumerOperator_ListHealthStatus_Call) RunAndReturn(run func(context.Context, Identity) ([]HealthStatus, error)) *MockServiceConsumerOperator_ListHealthStatus_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockServiceConsumerOperator interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockServiceConsumerOperator creates a new instance of MockServiceConsumerOperator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockServiceConsumerOperator(t mockConstructorTestingTNewMockServiceConsumerOperator) *MockServiceConsumerOperator {
	mock := &MockServiceConsumerOperator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
