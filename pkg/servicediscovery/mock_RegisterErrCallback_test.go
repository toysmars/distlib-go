// Code generated by mockery v2.26.1. DO NOT EDIT.

package servicediscovery

import mock "github.com/stretchr/testify/mock"

// MockRegisterErrCallback is an autogenerated mock type for the RegisterErrCallback type
type MockRegisterErrCallback struct {
	mock.Mock
}

type MockRegisterErrCallback_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegisterErrCallback) EXPECT() *MockRegisterErrCallback_Expecter {
	return &MockRegisterErrCallback_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: err
func (_m *MockRegisterErrCallback) Execute(err error) {
	_m.Called(err)
}

// MockRegisterErrCallback_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockRegisterErrCallback_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - err error
func (_e *MockRegisterErrCallback_Expecter) Execute(err interface{}) *MockRegisterErrCallback_Execute_Call {
	return &MockRegisterErrCallback_Execute_Call{Call: _e.mock.On("Execute", err)}
}

func (_c *MockRegisterErrCallback_Execute_Call) Run(run func(err error)) *MockRegisterErrCallback_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *MockRegisterErrCallback_Execute_Call) Return() *MockRegisterErrCallback_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockRegisterErrCallback_Execute_Call) RunAndReturn(run func(error)) *MockRegisterErrCallback_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRegisterErrCallback interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRegisterErrCallback creates a new instance of MockRegisterErrCallback. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRegisterErrCallback(t mockConstructorTestingTNewMockRegisterErrCallback) *MockRegisterErrCallback {
	mock := &MockRegisterErrCallback{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}