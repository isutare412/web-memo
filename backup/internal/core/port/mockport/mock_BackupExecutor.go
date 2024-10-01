// Code generated by mockery v2.46.1. DO NOT EDIT.

package mockport

import (
	context "context"

	model "github.com/isutare412/web-memo/backup/internal/core/model"
	mock "github.com/stretchr/testify/mock"
)

// MockBackupExecutor is an autogenerated mock type for the BackupExecutor type
type MockBackupExecutor struct {
	mock.Mock
}

type MockBackupExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBackupExecutor) EXPECT() *MockBackupExecutor_Expecter {
	return &MockBackupExecutor_Expecter{mock: &_m.Mock}
}

// BackupDatabase provides a mock function with given fields: _a0, _a1
func (_m *MockBackupExecutor) BackupDatabase(_a0 context.Context, _a1 model.DatabaseBackupRequest) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BackupDatabase")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.DatabaseBackupRequest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBackupExecutor_BackupDatabase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BackupDatabase'
type MockBackupExecutor_BackupDatabase_Call struct {
	*mock.Call
}

// BackupDatabase is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.DatabaseBackupRequest
func (_e *MockBackupExecutor_Expecter) BackupDatabase(_a0 interface{}, _a1 interface{}) *MockBackupExecutor_BackupDatabase_Call {
	return &MockBackupExecutor_BackupDatabase_Call{Call: _e.mock.On("BackupDatabase", _a0, _a1)}
}

func (_c *MockBackupExecutor_BackupDatabase_Call) Run(run func(_a0 context.Context, _a1 model.DatabaseBackupRequest)) *MockBackupExecutor_BackupDatabase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.DatabaseBackupRequest))
	})
	return _c
}

func (_c *MockBackupExecutor_BackupDatabase_Call) Return(_a0 error) *MockBackupExecutor_BackupDatabase_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBackupExecutor_BackupDatabase_Call) RunAndReturn(run func(context.Context, model.DatabaseBackupRequest) error) *MockBackupExecutor_BackupDatabase_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBackupExecutor creates a new instance of MockBackupExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBackupExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBackupExecutor {
	mock := &MockBackupExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
