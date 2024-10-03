// Code generated by mockery v2.46.1. DO NOT EDIT.

package mockport

import (
	context "context"

	ent "github.com/isutare412/web-memo/api/internal/core/ent"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockTagRepository is an autogenerated mock type for the TagRepository type
type MockTagRepository struct {
	mock.Mock
}

type MockTagRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTagRepository) EXPECT() *MockTagRepository_Expecter {
	return &MockTagRepository_Expecter{mock: &_m.Mock}
}

// CreateIfNotExist provides a mock function with given fields: ctx, tagName
func (_m *MockTagRepository) CreateIfNotExist(ctx context.Context, tagName string) (*ent.Tag, error) {
	ret := _m.Called(ctx, tagName)

	if len(ret) == 0 {
		panic("no return value specified for CreateIfNotExist")
	}

	var r0 *ent.Tag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ent.Tag, error)); ok {
		return rf(ctx, tagName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ent.Tag); ok {
		r0 = rf(ctx, tagName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Tag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tagName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTagRepository_CreateIfNotExist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateIfNotExist'
type MockTagRepository_CreateIfNotExist_Call struct {
	*mock.Call
}

// CreateIfNotExist is a helper method to define mock.On call
//   - ctx context.Context
//   - tagName string
func (_e *MockTagRepository_Expecter) CreateIfNotExist(ctx interface{}, tagName interface{}) *MockTagRepository_CreateIfNotExist_Call {
	return &MockTagRepository_CreateIfNotExist_Call{Call: _e.mock.On("CreateIfNotExist", ctx, tagName)}
}

func (_c *MockTagRepository_CreateIfNotExist_Call) Run(run func(ctx context.Context, tagName string)) *MockTagRepository_CreateIfNotExist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTagRepository_CreateIfNotExist_Call) Return(_a0 *ent.Tag, _a1 error) *MockTagRepository_CreateIfNotExist_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTagRepository_CreateIfNotExist_Call) RunAndReturn(run func(context.Context, string) (*ent.Tag, error)) *MockTagRepository_CreateIfNotExist_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllWithoutMemo provides a mock function with given fields: ctx, excludes
func (_m *MockTagRepository) DeleteAllWithoutMemo(ctx context.Context, excludes []string) (int, error) {
	ret := _m.Called(ctx, excludes)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllWithoutMemo")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) (int, error)); ok {
		return rf(ctx, excludes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) int); ok {
		r0 = rf(ctx, excludes)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, excludes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTagRepository_DeleteAllWithoutMemo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllWithoutMemo'
type MockTagRepository_DeleteAllWithoutMemo_Call struct {
	*mock.Call
}

// DeleteAllWithoutMemo is a helper method to define mock.On call
//   - ctx context.Context
//   - excludes []string
func (_e *MockTagRepository_Expecter) DeleteAllWithoutMemo(ctx interface{}, excludes interface{}) *MockTagRepository_DeleteAllWithoutMemo_Call {
	return &MockTagRepository_DeleteAllWithoutMemo_Call{Call: _e.mock.On("DeleteAllWithoutMemo", ctx, excludes)}
}

func (_c *MockTagRepository_DeleteAllWithoutMemo_Call) Run(run func(ctx context.Context, excludes []string)) *MockTagRepository_DeleteAllWithoutMemo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockTagRepository_DeleteAllWithoutMemo_Call) Return(count int, err error) *MockTagRepository_DeleteAllWithoutMemo_Call {
	_c.Call.Return(count, err)
	return _c
}

func (_c *MockTagRepository_DeleteAllWithoutMemo_Call) RunAndReturn(run func(context.Context, []string) (int, error)) *MockTagRepository_DeleteAllWithoutMemo_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllByMemoID provides a mock function with given fields: ctx, memoID
func (_m *MockTagRepository) FindAllByMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.Tag, error) {
	ret := _m.Called(ctx, memoID)

	if len(ret) == 0 {
		panic("no return value specified for FindAllByMemoID")
	}

	var r0 []*ent.Tag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*ent.Tag, error)); ok {
		return rf(ctx, memoID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*ent.Tag); ok {
		r0 = rf(ctx, memoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Tag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, memoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTagRepository_FindAllByMemoID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllByMemoID'
type MockTagRepository_FindAllByMemoID_Call struct {
	*mock.Call
}

// FindAllByMemoID is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
func (_e *MockTagRepository_Expecter) FindAllByMemoID(ctx interface{}, memoID interface{}) *MockTagRepository_FindAllByMemoID_Call {
	return &MockTagRepository_FindAllByMemoID_Call{Call: _e.mock.On("FindAllByMemoID", ctx, memoID)}
}

func (_c *MockTagRepository_FindAllByMemoID_Call) Run(run func(ctx context.Context, memoID uuid.UUID)) *MockTagRepository_FindAllByMemoID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockTagRepository_FindAllByMemoID_Call) Return(_a0 []*ent.Tag, _a1 error) *MockTagRepository_FindAllByMemoID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTagRepository_FindAllByMemoID_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]*ent.Tag, error)) *MockTagRepository_FindAllByMemoID_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllByUserIDAndNameContains provides a mock function with given fields: ctx, userID, name
func (_m *MockTagRepository) FindAllByUserIDAndNameContains(ctx context.Context, userID uuid.UUID, name string) ([]*ent.Tag, error) {
	ret := _m.Called(ctx, userID, name)

	if len(ret) == 0 {
		panic("no return value specified for FindAllByUserIDAndNameContains")
	}

	var r0 []*ent.Tag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) ([]*ent.Tag, error)); ok {
		return rf(ctx, userID, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) []*ent.Tag); ok {
		r0 = rf(ctx, userID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Tag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, userID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTagRepository_FindAllByUserIDAndNameContains_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllByUserIDAndNameContains'
type MockTagRepository_FindAllByUserIDAndNameContains_Call struct {
	*mock.Call
}

// FindAllByUserIDAndNameContains is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - name string
func (_e *MockTagRepository_Expecter) FindAllByUserIDAndNameContains(ctx interface{}, userID interface{}, name interface{}) *MockTagRepository_FindAllByUserIDAndNameContains_Call {
	return &MockTagRepository_FindAllByUserIDAndNameContains_Call{Call: _e.mock.On("FindAllByUserIDAndNameContains", ctx, userID, name)}
}

func (_c *MockTagRepository_FindAllByUserIDAndNameContains_Call) Run(run func(ctx context.Context, userID uuid.UUID, name string)) *MockTagRepository_FindAllByUserIDAndNameContains_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *MockTagRepository_FindAllByUserIDAndNameContains_Call) Return(_a0 []*ent.Tag, _a1 error) *MockTagRepository_FindAllByUserIDAndNameContains_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTagRepository_FindAllByUserIDAndNameContains_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) ([]*ent.Tag, error)) *MockTagRepository_FindAllByUserIDAndNameContains_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTagRepository creates a new instance of MockTagRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTagRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTagRepository {
	mock := &MockTagRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}