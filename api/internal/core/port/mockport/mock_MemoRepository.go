// Code generated by mockery v2.46.1. DO NOT EDIT.

package mockport

import (
	context "context"

	ent "github.com/isutare412/web-memo/api/internal/core/ent"
	mock "github.com/stretchr/testify/mock"

	model "github.com/isutare412/web-memo/api/internal/core/model"

	uuid "github.com/google/uuid"
)

// MockMemoRepository is an autogenerated mock type for the MemoRepository type
type MockMemoRepository struct {
	mock.Mock
}

type MockMemoRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMemoRepository) EXPECT() *MockMemoRepository_Expecter {
	return &MockMemoRepository_Expecter{mock: &_m.Mock}
}

// ClearSubscribers provides a mock function with given fields: ctx, memoID
func (_m *MockMemoRepository) ClearSubscribers(ctx context.Context, memoID uuid.UUID) error {
	ret := _m.Called(ctx, memoID)

	if len(ret) == 0 {
		panic("no return value specified for ClearSubscribers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, memoID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMemoRepository_ClearSubscribers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearSubscribers'
type MockMemoRepository_ClearSubscribers_Call struct {
	*mock.Call
}

// ClearSubscribers is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
func (_e *MockMemoRepository_Expecter) ClearSubscribers(ctx interface{}, memoID interface{}) *MockMemoRepository_ClearSubscribers_Call {
	return &MockMemoRepository_ClearSubscribers_Call{Call: _e.mock.On("ClearSubscribers", ctx, memoID)}
}

func (_c *MockMemoRepository_ClearSubscribers_Call) Run(run func(ctx context.Context, memoID uuid.UUID)) *MockMemoRepository_ClearSubscribers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_ClearSubscribers_Call) Return(_a0 error) *MockMemoRepository_ClearSubscribers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMemoRepository_ClearSubscribers_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *MockMemoRepository_ClearSubscribers_Call {
	_c.Call.Return(run)
	return _c
}

// CountByUserIDAndTagNames provides a mock function with given fields: ctx, userID, tags
func (_m *MockMemoRepository) CountByUserIDAndTagNames(ctx context.Context, userID uuid.UUID, tags []string) (int, error) {
	ret := _m.Called(ctx, userID, tags)

	if len(ret) == 0 {
		panic("no return value specified for CountByUserIDAndTagNames")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []string) (int, error)); ok {
		return rf(ctx, userID, tags)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []string) int); ok {
		r0 = rf(ctx, userID, tags)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, []string) error); ok {
		r1 = rf(ctx, userID, tags)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_CountByUserIDAndTagNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountByUserIDAndTagNames'
type MockMemoRepository_CountByUserIDAndTagNames_Call struct {
	*mock.Call
}

// CountByUserIDAndTagNames is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - tags []string
func (_e *MockMemoRepository_Expecter) CountByUserIDAndTagNames(ctx interface{}, userID interface{}, tags interface{}) *MockMemoRepository_CountByUserIDAndTagNames_Call {
	return &MockMemoRepository_CountByUserIDAndTagNames_Call{Call: _e.mock.On("CountByUserIDAndTagNames", ctx, userID, tags)}
}

func (_c *MockMemoRepository_CountByUserIDAndTagNames_Call) Run(run func(ctx context.Context, userID uuid.UUID, tags []string)) *MockMemoRepository_CountByUserIDAndTagNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].([]string))
	})
	return _c
}

func (_c *MockMemoRepository_CountByUserIDAndTagNames_Call) Return(_a0 int, _a1 error) *MockMemoRepository_CountByUserIDAndTagNames_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_CountByUserIDAndTagNames_Call) RunAndReturn(run func(context.Context, uuid.UUID, []string) (int, error)) *MockMemoRepository_CountByUserIDAndTagNames_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, memo, userID, tagIDs
func (_m *MockMemoRepository) Create(ctx context.Context, memo *ent.Memo, userID uuid.UUID, tagIDs []int) (*ent.Memo, error) {
	ret := _m.Called(ctx, memo, userID, tagIDs)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ent.Memo, uuid.UUID, []int) (*ent.Memo, error)); ok {
		return rf(ctx, memo, userID, tagIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ent.Memo, uuid.UUID, []int) *ent.Memo); ok {
		r0 = rf(ctx, memo, userID, tagIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ent.Memo, uuid.UUID, []int) error); ok {
		r1 = rf(ctx, memo, userID, tagIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockMemoRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - memo *ent.Memo
//   - userID uuid.UUID
//   - tagIDs []int
func (_e *MockMemoRepository_Expecter) Create(ctx interface{}, memo interface{}, userID interface{}, tagIDs interface{}) *MockMemoRepository_Create_Call {
	return &MockMemoRepository_Create_Call{Call: _e.mock.On("Create", ctx, memo, userID, tagIDs)}
}

func (_c *MockMemoRepository_Create_Call) Run(run func(ctx context.Context, memo *ent.Memo, userID uuid.UUID, tagIDs []int)) *MockMemoRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ent.Memo), args[2].(uuid.UUID), args[3].([]int))
	})
	return _c
}

func (_c *MockMemoRepository_Create_Call) Return(_a0 *ent.Memo, _a1 error) *MockMemoRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_Create_Call) RunAndReturn(run func(context.Context, *ent.Memo, uuid.UUID, []int) (*ent.Memo, error)) *MockMemoRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, memoID
func (_m *MockMemoRepository) Delete(ctx context.Context, memoID uuid.UUID) error {
	ret := _m.Called(ctx, memoID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, memoID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMemoRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockMemoRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
func (_e *MockMemoRepository_Expecter) Delete(ctx interface{}, memoID interface{}) *MockMemoRepository_Delete_Call {
	return &MockMemoRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, memoID)}
}

func (_c *MockMemoRepository_Delete_Call) Run(run func(ctx context.Context, memoID uuid.UUID)) *MockMemoRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_Delete_Call) Return(_a0 error) *MockMemoRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMemoRepository_Delete_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *MockMemoRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllByUserIDAndTagNamesWithEdges provides a mock function with given fields: ctx, userID, tags, sortParams, pageParams
func (_m *MockMemoRepository) FindAllByUserIDAndTagNamesWithEdges(ctx context.Context, userID uuid.UUID, tags []string, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error) {
	ret := _m.Called(ctx, userID, tags, sortParams, pageParams)

	if len(ret) == 0 {
		panic("no return value specified for FindAllByUserIDAndTagNamesWithEdges")
	}

	var r0 []*ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []string, model.MemoSortParams, model.PaginationParams) ([]*ent.Memo, error)); ok {
		return rf(ctx, userID, tags, sortParams, pageParams)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []string, model.MemoSortParams, model.PaginationParams) []*ent.Memo); ok {
		r0 = rf(ctx, userID, tags, sortParams, pageParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, []string, model.MemoSortParams, model.PaginationParams) error); ok {
		r1 = rf(ctx, userID, tags, sortParams, pageParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllByUserIDAndTagNamesWithEdges'
type MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call struct {
	*mock.Call
}

// FindAllByUserIDAndTagNamesWithEdges is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - tags []string
//   - sortParams model.MemoSortParams
//   - pageParams model.PaginationParams
func (_e *MockMemoRepository_Expecter) FindAllByUserIDAndTagNamesWithEdges(ctx interface{}, userID interface{}, tags interface{}, sortParams interface{}, pageParams interface{}) *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call {
	return &MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call{Call: _e.mock.On("FindAllByUserIDAndTagNamesWithEdges", ctx, userID, tags, sortParams, pageParams)}
}

func (_c *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call) Run(run func(ctx context.Context, userID uuid.UUID, tags []string, sortParams model.MemoSortParams, pageParams model.PaginationParams)) *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].([]string), args[3].(model.MemoSortParams), args[4].(model.PaginationParams))
	})
	return _c
}

func (_c *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call) Return(_a0 []*ent.Memo, _a1 error) *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call) RunAndReturn(run func(context.Context, uuid.UUID, []string, model.MemoSortParams, model.PaginationParams) ([]*ent.Memo, error)) *MockMemoRepository_FindAllByUserIDAndTagNamesWithEdges_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllByUserIDWithEdges provides a mock function with given fields: ctx, userID, sortParams, pageParams
func (_m *MockMemoRepository) FindAllByUserIDWithEdges(ctx context.Context, userID uuid.UUID, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error) {
	ret := _m.Called(ctx, userID, sortParams, pageParams)

	if len(ret) == 0 {
		panic("no return value specified for FindAllByUserIDWithEdges")
	}

	var r0 []*ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.MemoSortParams, model.PaginationParams) ([]*ent.Memo, error)); ok {
		return rf(ctx, userID, sortParams, pageParams)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.MemoSortParams, model.PaginationParams) []*ent.Memo); ok {
		r0 = rf(ctx, userID, sortParams, pageParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, model.MemoSortParams, model.PaginationParams) error); ok {
		r1 = rf(ctx, userID, sortParams, pageParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_FindAllByUserIDWithEdges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllByUserIDWithEdges'
type MockMemoRepository_FindAllByUserIDWithEdges_Call struct {
	*mock.Call
}

// FindAllByUserIDWithEdges is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - sortParams model.MemoSortParams
//   - pageParams model.PaginationParams
func (_e *MockMemoRepository_Expecter) FindAllByUserIDWithEdges(ctx interface{}, userID interface{}, sortParams interface{}, pageParams interface{}) *MockMemoRepository_FindAllByUserIDWithEdges_Call {
	return &MockMemoRepository_FindAllByUserIDWithEdges_Call{Call: _e.mock.On("FindAllByUserIDWithEdges", ctx, userID, sortParams, pageParams)}
}

func (_c *MockMemoRepository_FindAllByUserIDWithEdges_Call) Run(run func(ctx context.Context, userID uuid.UUID, sortParams model.MemoSortParams, pageParams model.PaginationParams)) *MockMemoRepository_FindAllByUserIDWithEdges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(model.MemoSortParams), args[3].(model.PaginationParams))
	})
	return _c
}

func (_c *MockMemoRepository_FindAllByUserIDWithEdges_Call) Return(_a0 []*ent.Memo, _a1 error) *MockMemoRepository_FindAllByUserIDWithEdges_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_FindAllByUserIDWithEdges_Call) RunAndReturn(run func(context.Context, uuid.UUID, model.MemoSortParams, model.PaginationParams) ([]*ent.Memo, error)) *MockMemoRepository_FindAllByUserIDWithEdges_Call {
	_c.Call.Return(run)
	return _c
}

// FindByID provides a mock function with given fields: ctx, memoID
func (_m *MockMemoRepository) FindByID(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	ret := _m.Called(ctx, memoID)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*ent.Memo, error)); ok {
		return rf(ctx, memoID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *ent.Memo); ok {
		r0 = rf(ctx, memoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, memoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type MockMemoRepository_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
func (_e *MockMemoRepository_Expecter) FindByID(ctx interface{}, memoID interface{}) *MockMemoRepository_FindByID_Call {
	return &MockMemoRepository_FindByID_Call{Call: _e.mock.On("FindByID", ctx, memoID)}
}

func (_c *MockMemoRepository_FindByID_Call) Run(run func(ctx context.Context, memoID uuid.UUID)) *MockMemoRepository_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_FindByID_Call) Return(_a0 *ent.Memo, _a1 error) *MockMemoRepository_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_FindByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*ent.Memo, error)) *MockMemoRepository_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// FindByIDWithEdges provides a mock function with given fields: ctx, memoID
func (_m *MockMemoRepository) FindByIDWithEdges(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	ret := _m.Called(ctx, memoID)

	if len(ret) == 0 {
		panic("no return value specified for FindByIDWithEdges")
	}

	var r0 *ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*ent.Memo, error)); ok {
		return rf(ctx, memoID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *ent.Memo); ok {
		r0 = rf(ctx, memoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, memoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_FindByIDWithEdges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByIDWithEdges'
type MockMemoRepository_FindByIDWithEdges_Call struct {
	*mock.Call
}

// FindByIDWithEdges is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
func (_e *MockMemoRepository_Expecter) FindByIDWithEdges(ctx interface{}, memoID interface{}) *MockMemoRepository_FindByIDWithEdges_Call {
	return &MockMemoRepository_FindByIDWithEdges_Call{Call: _e.mock.On("FindByIDWithEdges", ctx, memoID)}
}

func (_c *MockMemoRepository_FindByIDWithEdges_Call) Run(run func(ctx context.Context, memoID uuid.UUID)) *MockMemoRepository_FindByIDWithEdges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_FindByIDWithEdges_Call) Return(_a0 *ent.Memo, _a1 error) *MockMemoRepository_FindByIDWithEdges_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_FindByIDWithEdges_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*ent.Memo, error)) *MockMemoRepository_FindByIDWithEdges_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterSubscriber provides a mock function with given fields: ctx, memoID, userID
func (_m *MockMemoRepository) RegisterSubscriber(ctx context.Context, memoID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(ctx, memoID, userID)

	if len(ret) == 0 {
		panic("no return value specified for RegisterSubscriber")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, memoID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMemoRepository_RegisterSubscriber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterSubscriber'
type MockMemoRepository_RegisterSubscriber_Call struct {
	*mock.Call
}

// RegisterSubscriber is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
//   - userID uuid.UUID
func (_e *MockMemoRepository_Expecter) RegisterSubscriber(ctx interface{}, memoID interface{}, userID interface{}) *MockMemoRepository_RegisterSubscriber_Call {
	return &MockMemoRepository_RegisterSubscriber_Call{Call: _e.mock.On("RegisterSubscriber", ctx, memoID, userID)}
}

func (_c *MockMemoRepository_RegisterSubscriber_Call) Run(run func(ctx context.Context, memoID uuid.UUID, userID uuid.UUID)) *MockMemoRepository_RegisterSubscriber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_RegisterSubscriber_Call) Return(_a0 error) *MockMemoRepository_RegisterSubscriber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMemoRepository_RegisterSubscriber_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MockMemoRepository_RegisterSubscriber_Call {
	_c.Call.Return(run)
	return _c
}

// ReplaceTags provides a mock function with given fields: ctx, memoID, tagIDs, updateTime
func (_m *MockMemoRepository) ReplaceTags(ctx context.Context, memoID uuid.UUID, tagIDs []int, updateTime bool) error {
	ret := _m.Called(ctx, memoID, tagIDs, updateTime)

	if len(ret) == 0 {
		panic("no return value specified for ReplaceTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []int, bool) error); ok {
		r0 = rf(ctx, memoID, tagIDs, updateTime)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMemoRepository_ReplaceTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReplaceTags'
type MockMemoRepository_ReplaceTags_Call struct {
	*mock.Call
}

// ReplaceTags is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
//   - tagIDs []int
//   - updateTime bool
func (_e *MockMemoRepository_Expecter) ReplaceTags(ctx interface{}, memoID interface{}, tagIDs interface{}, updateTime interface{}) *MockMemoRepository_ReplaceTags_Call {
	return &MockMemoRepository_ReplaceTags_Call{Call: _e.mock.On("ReplaceTags", ctx, memoID, tagIDs, updateTime)}
}

func (_c *MockMemoRepository_ReplaceTags_Call) Run(run func(ctx context.Context, memoID uuid.UUID, tagIDs []int, updateTime bool)) *MockMemoRepository_ReplaceTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].([]int), args[3].(bool))
	})
	return _c
}

func (_c *MockMemoRepository_ReplaceTags_Call) Return(_a0 error) *MockMemoRepository_ReplaceTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMemoRepository_ReplaceTags_Call) RunAndReturn(run func(context.Context, uuid.UUID, []int, bool) error) *MockMemoRepository_ReplaceTags_Call {
	_c.Call.Return(run)
	return _c
}

// UnregisterSubscriber provides a mock function with given fields: ctx, memoID, userID
func (_m *MockMemoRepository) UnregisterSubscriber(ctx context.Context, memoID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(ctx, memoID, userID)

	if len(ret) == 0 {
		panic("no return value specified for UnregisterSubscriber")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, memoID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMemoRepository_UnregisterSubscriber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnregisterSubscriber'
type MockMemoRepository_UnregisterSubscriber_Call struct {
	*mock.Call
}

// UnregisterSubscriber is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
//   - userID uuid.UUID
func (_e *MockMemoRepository_Expecter) UnregisterSubscriber(ctx interface{}, memoID interface{}, userID interface{}) *MockMemoRepository_UnregisterSubscriber_Call {
	return &MockMemoRepository_UnregisterSubscriber_Call{Call: _e.mock.On("UnregisterSubscriber", ctx, memoID, userID)}
}

func (_c *MockMemoRepository_UnregisterSubscriber_Call) Run(run func(ctx context.Context, memoID uuid.UUID, userID uuid.UUID)) *MockMemoRepository_UnregisterSubscriber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockMemoRepository_UnregisterSubscriber_Call) Return(_a0 error) *MockMemoRepository_UnregisterSubscriber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMemoRepository_UnregisterSubscriber_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MockMemoRepository_UnregisterSubscriber_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *MockMemoRepository) Update(_a0 context.Context, _a1 *ent.Memo) (*ent.Memo, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ent.Memo) (*ent.Memo, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ent.Memo) *ent.Memo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ent.Memo) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockMemoRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ent.Memo
func (_e *MockMemoRepository_Expecter) Update(_a0 interface{}, _a1 interface{}) *MockMemoRepository_Update_Call {
	return &MockMemoRepository_Update_Call{Call: _e.mock.On("Update", _a0, _a1)}
}

func (_c *MockMemoRepository_Update_Call) Run(run func(_a0 context.Context, _a1 *ent.Memo)) *MockMemoRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ent.Memo))
	})
	return _c
}

func (_c *MockMemoRepository_Update_Call) Return(_a0 *ent.Memo, _a1 error) *MockMemoRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_Update_Call) RunAndReturn(run func(context.Context, *ent.Memo) (*ent.Memo, error)) *MockMemoRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateIsPublish provides a mock function with given fields: ctx, memoID, isPublish
func (_m *MockMemoRepository) UpdateIsPublish(ctx context.Context, memoID uuid.UUID, isPublish bool) (*ent.Memo, error) {
	ret := _m.Called(ctx, memoID, isPublish)

	if len(ret) == 0 {
		panic("no return value specified for UpdateIsPublish")
	}

	var r0 *ent.Memo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) (*ent.Memo, error)); ok {
		return rf(ctx, memoID, isPublish)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) *ent.Memo); ok {
		r0 = rf(ctx, memoID, isPublish)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Memo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, bool) error); ok {
		r1 = rf(ctx, memoID, isPublish)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMemoRepository_UpdateIsPublish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateIsPublish'
type MockMemoRepository_UpdateIsPublish_Call struct {
	*mock.Call
}

// UpdateIsPublish is a helper method to define mock.On call
//   - ctx context.Context
//   - memoID uuid.UUID
//   - isPublish bool
func (_e *MockMemoRepository_Expecter) UpdateIsPublish(ctx interface{}, memoID interface{}, isPublish interface{}) *MockMemoRepository_UpdateIsPublish_Call {
	return &MockMemoRepository_UpdateIsPublish_Call{Call: _e.mock.On("UpdateIsPublish", ctx, memoID, isPublish)}
}

func (_c *MockMemoRepository_UpdateIsPublish_Call) Run(run func(ctx context.Context, memoID uuid.UUID, isPublish bool)) *MockMemoRepository_UpdateIsPublish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(bool))
	})
	return _c
}

func (_c *MockMemoRepository_UpdateIsPublish_Call) Return(_a0 *ent.Memo, _a1 error) *MockMemoRepository_UpdateIsPublish_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMemoRepository_UpdateIsPublish_Call) RunAndReturn(run func(context.Context, uuid.UUID, bool) (*ent.Memo, error)) *MockMemoRepository_UpdateIsPublish_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMemoRepository creates a new instance of MockMemoRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMemoRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMemoRepository {
	mock := &MockMemoRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}