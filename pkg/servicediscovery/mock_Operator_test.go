// Code generated by mockery v2.26.1. DO NOT EDIT.

package servicediscovery

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockOperator is an autogenerated mock type for the Operator type
type MockOperator struct {
	mock.Mock
}

type MockOperator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOperator) EXPECT() *MockOperator_Expecter {
	return &MockOperator_Expecter{mock: &_m.Mock}
}

// DeregisterHealthStatus provides a mock function with given fields: ctx, identity
func (_m *MockOperator) DeregisterHealthStatus(ctx context.Context, identity Identity) error {
	ret := _m.Called(ctx, identity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Identity) error); ok {
		r0 = rf(ctx, identity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOperator_DeregisterHealthStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeregisterHealthStatus'
type MockOperator_DeregisterHealthStatus_Call struct {
	*mock.Call
}

// DeregisterHealthStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - identity Identity
func (_e *MockOperator_Expecter) DeregisterHealthStatus(ctx interface{}, identity interface{}) *MockOperator_DeregisterHealthStatus_Call {
	return &MockOperator_DeregisterHealthStatus_Call{Call: _e.mock.On("DeregisterHealthStatus", ctx, identity)}
}

func (_c *MockOperator_DeregisterHealthStatus_Call) Run(run func(ctx context.Context, identity Identity)) *MockOperator_DeregisterHealthStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Identity))
	})
	return _c
}

func (_c *MockOperator_DeregisterHealthStatus_Call) Return(_a0 error) *MockOperator_DeregisterHealthStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOperator_DeregisterHealthStatus_Call) RunAndReturn(run func(context.Context, Identity) error) *MockOperator_DeregisterHealthStatus_Call {
	_c.Call.Return(run)
	return _c
}

// ListHealthStatus provides a mock function with given fields: ctx, identity
func (_m *MockOperator) ListHealthStatus(ctx context.Context, identity Identity) ([]HealthStatus, error) {
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

// MockOperator_ListHealthStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListHealthStatus'
type MockOperator_ListHealthStatus_Call struct {
	*mock.Call
}

// ListHealthStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - identity Identity
func (_e *MockOperator_Expecter) ListHealthStatus(ctx interface{}, identity interface{}) *MockOperator_ListHealthStatus_Call {
	return &MockOperator_ListHealthStatus_Call{Call: _e.mock.On("ListHealthStatus", ctx, identity)}
}

func (_c *MockOperator_ListHealthStatus_Call) Run(run func(ctx context.Context, identity Identity)) *MockOperator_ListHealthStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Identity))
	})
	return _c
}

func (_c *MockOperator_ListHealthStatus_Call) Return(_a0 []HealthStatus, _a1 error) *MockOperator_ListHealthStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOperator_ListHealthStatus_Call) RunAndReturn(run func(context.Context, Identity) ([]HealthStatus, error)) *MockOperator_ListHealthStatus_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterHealthStatus provides a mock function with given fields: ctx, status
func (_m *MockOperator) RegisterHealthStatus(ctx context.Context, status HealthStatus) error {
	ret := _m.Called(ctx, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, HealthStatus) error); ok {
		r0 = rf(ctx, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOperator_RegisterHealthStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterHealthStatus'
type MockOperator_RegisterHealthStatus_Call struct {
	*mock.Call
}

// RegisterHealthStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - status HealthStatus
func (_e *MockOperator_Expecter) RegisterHealthStatus(ctx interface{}, status interface{}) *MockOperator_RegisterHealthStatus_Call {
	return &MockOperator_RegisterHealthStatus_Call{Call: _e.mock.On("RegisterHealthStatus", ctx, status)}
}

func (_c *MockOperator_RegisterHealthStatus_Call) Run(run func(ctx context.Context, status HealthStatus)) *MockOperator_RegisterHealthStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(HealthStatus))
	})
	return _c
}

func (_c *MockOperator_RegisterHealthStatus_Call) Return(_a0 error) *MockOperator_RegisterHealthStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOperator_RegisterHealthStatus_Call) RunAndReturn(run func(context.Context, HealthStatus) error) *MockOperator_RegisterHealthStatus_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockOperator interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockOperator creates a new instance of MockOperator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockOperator(t mockConstructorTestingTNewMockOperator) *MockOperator {
	mock := &MockOperator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
