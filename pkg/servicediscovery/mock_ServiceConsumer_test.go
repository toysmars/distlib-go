// Code generated by mockery v2.26.1. DO NOT EDIT.

package servicediscovery

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockServiceConsumer is an autogenerated mock type for the ServiceConsumer type
type MockServiceConsumer struct {
	mock.Mock
}

type MockServiceConsumer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockServiceConsumer) EXPECT() *MockServiceConsumer_Expecter {
	return &MockServiceConsumer_Expecter{mock: &_m.Mock}
}

// List provides a mock function with given fields: ctx
func (_m *MockServiceConsumer) List(ctx context.Context) ([]HealthStatus, error) {
	ret := _m.Called(ctx)

	var r0 []HealthStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]HealthStatus, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []HealthStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]HealthStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceConsumer_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockServiceConsumer_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockServiceConsumer_Expecter) List(ctx interface{}) *MockServiceConsumer_List_Call {
	return &MockServiceConsumer_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *MockServiceConsumer_List_Call) Run(run func(ctx context.Context)) *MockServiceConsumer_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockServiceConsumer_List_Call) Return(_a0 []HealthStatus, _a1 error) *MockServiceConsumer_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceConsumer_List_Call) RunAndReturn(run func(context.Context) ([]HealthStatus, error)) *MockServiceConsumer_List_Call {
	_c.Call.Return(run)
	return _c
}

// Run provides a mock function with given fields: ctx
func (_m *MockServiceConsumer) Run(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServiceConsumer_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockServiceConsumer_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockServiceConsumer_Expecter) Run(ctx interface{}) *MockServiceConsumer_Run_Call {
	return &MockServiceConsumer_Run_Call{Call: _e.mock.On("Run", ctx)}
}

func (_c *MockServiceConsumer_Run_Call) Run(run func(ctx context.Context)) *MockServiceConsumer_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockServiceConsumer_Run_Call) Return(_a0 error) *MockServiceConsumer_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceConsumer_Run_Call) RunAndReturn(run func(context.Context) error) *MockServiceConsumer_Run_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockServiceConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockServiceConsumer creates a new instance of MockServiceConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockServiceConsumer(t mockConstructorTestingTNewMockServiceConsumer) *MockServiceConsumer {
	mock := &MockServiceConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
