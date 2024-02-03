// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/port/repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/port/repository.go -destination=internal/core/port/mockport/mock_repository.go -package=mockport
//
// Package mockport is a generated GoMock package.
package mockport

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	ent "github.com/isutare412/web-memo/api/internal/core/ent"
	model "github.com/isutare412/web-memo/api/internal/core/model"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactionManager is a mock of TransactionManager interface.
type MockTransactionManager struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionManagerMockRecorder
}

// MockTransactionManagerMockRecorder is the mock recorder for MockTransactionManager.
type MockTransactionManagerMockRecorder struct {
	mock *MockTransactionManager
}

// NewMockTransactionManager creates a new mock instance.
func NewMockTransactionManager(ctrl *gomock.Controller) *MockTransactionManager {
	mock := &MockTransactionManager{ctrl: ctrl}
	mock.recorder = &MockTransactionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionManager) EXPECT() *MockTransactionManagerMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *MockTransactionManager) BeginTx(arg0 context.Context) (context.Context, func() error, func() error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", arg0)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(func() error)
	ret2, _ := ret[2].(func() error)
	return ret0, ret1, ret2
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockTransactionManagerMockRecorder) BeginTx(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockTransactionManager)(nil).BeginTx), arg0)
}

// WithTx mocks base method.
func (m *MockTransactionManager) WithTx(arg0 context.Context, arg1 func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockTransactionManagerMockRecorder) WithTx(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockTransactionManager)(nil).WithTx), arg0, arg1)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*ent.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// FindByID mocks base method.
func (m *MockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*ent.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, userID)
	ret0, _ := ret[0].(*ent.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryMockRecorder) FindByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), ctx, userID)
}

// Upsert mocks base method.
func (m *MockUserRepository) Upsert(arg0 context.Context, arg1 *ent.User) (*ent.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1)
	ret0, _ := ret[0].(*ent.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upsert indicates an expected call of Upsert.
func (mr *MockUserRepositoryMockRecorder) Upsert(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockUserRepository)(nil).Upsert), arg0, arg1)
}

// MockMemoRepository is a mock of MemoRepository interface.
type MockMemoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMemoRepositoryMockRecorder
}

// MockMemoRepositoryMockRecorder is the mock recorder for MockMemoRepository.
type MockMemoRepositoryMockRecorder struct {
	mock *MockMemoRepository
}

// NewMockMemoRepository creates a new mock instance.
func NewMockMemoRepository(ctrl *gomock.Controller) *MockMemoRepository {
	mock := &MockMemoRepository{ctrl: ctrl}
	mock.recorder = &MockMemoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemoRepository) EXPECT() *MockMemoRepositoryMockRecorder {
	return m.recorder
}

// CountByUserIDAndTagNames mocks base method.
func (m *MockMemoRepository) CountByUserIDAndTagNames(ctx context.Context, userID uuid.UUID, tags []string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByUserIDAndTagNames", ctx, userID, tags)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByUserIDAndTagNames indicates an expected call of CountByUserIDAndTagNames.
func (mr *MockMemoRepositoryMockRecorder) CountByUserIDAndTagNames(ctx, userID, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByUserIDAndTagNames", reflect.TypeOf((*MockMemoRepository)(nil).CountByUserIDAndTagNames), ctx, userID, tags)
}

// Create mocks base method.
func (m *MockMemoRepository) Create(ctx context.Context, memo *ent.Memo, userID uuid.UUID, tagIDs []int) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, memo, userID, tagIDs)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockMemoRepositoryMockRecorder) Create(ctx, memo, userID, tagIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMemoRepository)(nil).Create), ctx, memo, userID, tagIDs)
}

// Delete mocks base method.
func (m *MockMemoRepository) Delete(ctx context.Context, memoID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, memoID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockMemoRepositoryMockRecorder) Delete(ctx, memoID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMemoRepository)(nil).Delete), ctx, memoID)
}

// FindAllByUserIDAndTagNamesWithTags mocks base method.
func (m *MockMemoRepository) FindAllByUserIDAndTagNamesWithTags(ctx context.Context, userID uuid.UUID, tags []string, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserIDAndTagNamesWithTags", ctx, userID, tags, sortParams, pageParams)
	ret0, _ := ret[0].([]*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserIDAndTagNamesWithTags indicates an expected call of FindAllByUserIDAndTagNamesWithTags.
func (mr *MockMemoRepositoryMockRecorder) FindAllByUserIDAndTagNamesWithTags(ctx, userID, tags, sortParams, pageParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserIDAndTagNamesWithTags", reflect.TypeOf((*MockMemoRepository)(nil).FindAllByUserIDAndTagNamesWithTags), ctx, userID, tags, sortParams, pageParams)
}

// FindAllByUserIDWithTags mocks base method.
func (m *MockMemoRepository) FindAllByUserIDWithTags(ctx context.Context, userID uuid.UUID, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserIDWithTags", ctx, userID, sortParams, pageParams)
	ret0, _ := ret[0].([]*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserIDWithTags indicates an expected call of FindAllByUserIDWithTags.
func (mr *MockMemoRepositoryMockRecorder) FindAllByUserIDWithTags(ctx, userID, sortParams, pageParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserIDWithTags", reflect.TypeOf((*MockMemoRepository)(nil).FindAllByUserIDWithTags), ctx, userID, sortParams, pageParams)
}

// FindByID mocks base method.
func (m *MockMemoRepository) FindByID(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, memoID)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockMemoRepositoryMockRecorder) FindByID(ctx, memoID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockMemoRepository)(nil).FindByID), ctx, memoID)
}

// FindByIDWithTags mocks base method.
func (m *MockMemoRepository) FindByIDWithTags(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDWithTags", ctx, memoID)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDWithTags indicates an expected call of FindByIDWithTags.
func (mr *MockMemoRepositoryMockRecorder) FindByIDWithTags(ctx, memoID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDWithTags", reflect.TypeOf((*MockMemoRepository)(nil).FindByIDWithTags), ctx, memoID)
}

// ReplaceTags mocks base method.
func (m *MockMemoRepository) ReplaceTags(ctx context.Context, memoID uuid.UUID, tagIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceTags", ctx, memoID, tagIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceTags indicates an expected call of ReplaceTags.
func (mr *MockMemoRepositoryMockRecorder) ReplaceTags(ctx, memoID, tagIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceTags", reflect.TypeOf((*MockMemoRepository)(nil).ReplaceTags), ctx, memoID, tagIDs)
}

// Update mocks base method.
func (m *MockMemoRepository) Update(arg0 context.Context, arg1 *ent.Memo) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockMemoRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMemoRepository)(nil).Update), arg0, arg1)
}

// MockTagRepository is a mock of TagRepository interface.
type MockTagRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTagRepositoryMockRecorder
}

// MockTagRepositoryMockRecorder is the mock recorder for MockTagRepository.
type MockTagRepositoryMockRecorder struct {
	mock *MockTagRepository
}

// NewMockTagRepository creates a new mock instance.
func NewMockTagRepository(ctrl *gomock.Controller) *MockTagRepository {
	mock := &MockTagRepository{ctrl: ctrl}
	mock.recorder = &MockTagRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagRepository) EXPECT() *MockTagRepositoryMockRecorder {
	return m.recorder
}

// CreateIfNotExist mocks base method.
func (m *MockTagRepository) CreateIfNotExist(ctx context.Context, tagName string) (*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIfNotExist", ctx, tagName)
	ret0, _ := ret[0].(*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIfNotExist indicates an expected call of CreateIfNotExist.
func (mr *MockTagRepositoryMockRecorder) CreateIfNotExist(ctx, tagName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIfNotExist", reflect.TypeOf((*MockTagRepository)(nil).CreateIfNotExist), ctx, tagName)
}

// DeleteAllWithoutMemo mocks base method.
func (m *MockTagRepository) DeleteAllWithoutMemo(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllWithoutMemo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAllWithoutMemo indicates an expected call of DeleteAllWithoutMemo.
func (mr *MockTagRepositoryMockRecorder) DeleteAllWithoutMemo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllWithoutMemo", reflect.TypeOf((*MockTagRepository)(nil).DeleteAllWithoutMemo), arg0)
}

// FindAllByMemoID mocks base method.
func (m *MockTagRepository) FindAllByMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByMemoID", ctx, memoID)
	ret0, _ := ret[0].([]*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByMemoID indicates an expected call of FindAllByMemoID.
func (mr *MockTagRepositoryMockRecorder) FindAllByMemoID(ctx, memoID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByMemoID", reflect.TypeOf((*MockTagRepository)(nil).FindAllByMemoID), ctx, memoID)
}

// FindAllByUserIDAndNameContains mocks base method.
func (m *MockTagRepository) FindAllByUserIDAndNameContains(ctx context.Context, userID uuid.UUID, name string) ([]*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserIDAndNameContains", ctx, userID, name)
	ret0, _ := ret[0].([]*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserIDAndNameContains indicates an expected call of FindAllByUserIDAndNameContains.
func (mr *MockTagRepositoryMockRecorder) FindAllByUserIDAndNameContains(ctx, userID, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserIDAndNameContains", reflect.TypeOf((*MockTagRepository)(nil).FindAllByUserIDAndNameContains), ctx, userID, name)
}

// MockKVRepository is a mock of KVRepository interface.
type MockKVRepository struct {
	ctrl     *gomock.Controller
	recorder *MockKVRepositoryMockRecorder
}

// MockKVRepositoryMockRecorder is the mock recorder for MockKVRepository.
type MockKVRepositoryMockRecorder struct {
	mock *MockKVRepository
}

// NewMockKVRepository creates a new mock instance.
func NewMockKVRepository(ctrl *gomock.Controller) *MockKVRepository {
	mock := &MockKVRepository{ctrl: ctrl}
	mock.recorder = &MockKVRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKVRepository) EXPECT() *MockKVRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockKVRepository) Delete(ctx context.Context, keys ...string) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockKVRepositoryMockRecorder) Delete(ctx any, keys ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKVRepository)(nil).Delete), varargs...)
}

// Get mocks base method.
func (m *MockKVRepository) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockKVRepositoryMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKVRepository)(nil).Get), ctx, key)
}

// GetThenDelete mocks base method.
func (m *MockKVRepository) GetThenDelete(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetThenDelete", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThenDelete indicates an expected call of GetThenDelete.
func (mr *MockKVRepositoryMockRecorder) GetThenDelete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThenDelete", reflect.TypeOf((*MockKVRepository)(nil).GetThenDelete), ctx, key)
}

// Set mocks base method.
func (m *MockKVRepository) Set(ctx context.Context, key, val string, exp time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, val, exp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockKVRepositoryMockRecorder) Set(ctx, key, val, exp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockKVRepository)(nil).Set), ctx, key, val, exp)
}
