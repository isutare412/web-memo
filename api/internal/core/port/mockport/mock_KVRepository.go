// Code generated by mockery v2.46.1. DO NOT EDIT.

package mockport

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockKVRepository is an autogenerated mock type for the KVRepository type
type MockKVRepository struct {
	mock.Mock
}

type MockKVRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKVRepository) EXPECT() *MockKVRepository_Expecter {
	return &MockKVRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, keys
func (_m *MockKVRepository) Delete(ctx context.Context, keys ...string) (int64, error) {
	_va := make([]interface{}, len(keys))
	for _i := range keys {
		_va[_i] = keys[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) (int64, error)); ok {
		return rf(ctx, keys...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...string) int64); ok {
		r0 = rf(ctx, keys...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, keys...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKVRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockKVRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - keys ...string
func (_e *MockKVRepository_Expecter) Delete(ctx interface{}, keys ...interface{}) *MockKVRepository_Delete_Call {
	return &MockKVRepository_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx}, keys...)...)}
}

func (_c *MockKVRepository_Delete_Call) Run(run func(ctx context.Context, keys ...string)) *MockKVRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MockKVRepository_Delete_Call) Return(delCount int64, err error) *MockKVRepository_Delete_Call {
	_c.Call.Return(delCount, err)
	return _c
}

func (_c *MockKVRepository_Delete_Call) RunAndReturn(run func(context.Context, ...string) (int64, error)) *MockKVRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key
func (_m *MockKVRepository) Get(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKVRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockKVRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockKVRepository_Expecter) Get(ctx interface{}, key interface{}) *MockKVRepository_Get_Call {
	return &MockKVRepository_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *MockKVRepository_Get_Call) Run(run func(ctx context.Context, key string)) *MockKVRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockKVRepository_Get_Call) Return(_a0 string, _a1 error) *MockKVRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKVRepository_Get_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockKVRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetThenDelete provides a mock function with given fields: ctx, key
func (_m *MockKVRepository) GetThenDelete(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GetThenDelete")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKVRepository_GetThenDelete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetThenDelete'
type MockKVRepository_GetThenDelete_Call struct {
	*mock.Call
}

// GetThenDelete is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockKVRepository_Expecter) GetThenDelete(ctx interface{}, key interface{}) *MockKVRepository_GetThenDelete_Call {
	return &MockKVRepository_GetThenDelete_Call{Call: _e.mock.On("GetThenDelete", ctx, key)}
}

func (_c *MockKVRepository_GetThenDelete_Call) Run(run func(ctx context.Context, key string)) *MockKVRepository_GetThenDelete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockKVRepository_GetThenDelete_Call) Return(_a0 string, _a1 error) *MockKVRepository_GetThenDelete_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKVRepository_GetThenDelete_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockKVRepository_GetThenDelete_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, val, exp
func (_m *MockKVRepository) Set(ctx context.Context, key string, val string, exp time.Duration) error {
	ret := _m.Called(ctx, key, val, exp)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Duration) error); ok {
		r0 = rf(ctx, key, val, exp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKVRepository_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockKVRepository_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - val string
//   - exp time.Duration
func (_e *MockKVRepository_Expecter) Set(ctx interface{}, key interface{}, val interface{}, exp interface{}) *MockKVRepository_Set_Call {
	return &MockKVRepository_Set_Call{Call: _e.mock.On("Set", ctx, key, val, exp)}
}

func (_c *MockKVRepository_Set_Call) Run(run func(ctx context.Context, key string, val string, exp time.Duration)) *MockKVRepository_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockKVRepository_Set_Call) Return(_a0 error) *MockKVRepository_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKVRepository_Set_Call) RunAndReturn(run func(context.Context, string, string, time.Duration) error) *MockKVRepository_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKVRepository creates a new instance of MockKVRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKVRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKVRepository {
	mock := &MockKVRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
