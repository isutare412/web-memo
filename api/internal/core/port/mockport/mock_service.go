// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/port/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/port/service.go -destination=internal/core/port/mockport/mock_service.go -package=mockport
//
// Package mockport is a generated GoMock package.
package mockport

import (
	context "context"
	http "net/http"
	reflect "reflect"

	uuid "github.com/google/uuid"
	ent "github.com/isutare412/web-memo/api/internal/core/ent"
	model "github.com/isutare412/web-memo/api/internal/core/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// FinishGoogleSignIn mocks base method.
func (m *MockAuthService) FinishGoogleSignIn(arg0 context.Context, arg1 *http.Request) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishGoogleSignIn", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FinishGoogleSignIn indicates an expected call of FinishGoogleSignIn.
func (mr *MockAuthServiceMockRecorder) FinishGoogleSignIn(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishGoogleSignIn", reflect.TypeOf((*MockAuthService)(nil).FinishGoogleSignIn), arg0, arg1)
}

// StartGoogleSignIn mocks base method.
func (m *MockAuthService) StartGoogleSignIn(arg0 context.Context, arg1 *http.Request) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartGoogleSignIn", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartGoogleSignIn indicates an expected call of StartGoogleSignIn.
func (mr *MockAuthServiceMockRecorder) StartGoogleSignIn(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartGoogleSignIn", reflect.TypeOf((*MockAuthService)(nil).StartGoogleSignIn), arg0, arg1)
}

// VerifyAppIDTokenString mocks base method.
func (m *MockAuthService) VerifyAppIDTokenString(arg0 string) (*model.AppIDToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyAppIDTokenString", arg0)
	ret0, _ := ret[0].(*model.AppIDToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyAppIDTokenString indicates an expected call of VerifyAppIDTokenString.
func (mr *MockAuthServiceMockRecorder) VerifyAppIDTokenString(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyAppIDTokenString", reflect.TypeOf((*MockAuthService)(nil).VerifyAppIDTokenString), arg0)
}

// MockMemoService is a mock of MemoService interface.
type MockMemoService struct {
	ctrl     *gomock.Controller
	recorder *MockMemoServiceMockRecorder
}

// MockMemoServiceMockRecorder is the mock recorder for MockMemoService.
type MockMemoServiceMockRecorder struct {
	mock *MockMemoService
}

// NewMockMemoService creates a new mock instance.
func NewMockMemoService(ctrl *gomock.Controller) *MockMemoService {
	mock := &MockMemoService{ctrl: ctrl}
	mock.recorder = &MockMemoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemoService) EXPECT() *MockMemoServiceMockRecorder {
	return m.recorder
}

// CreateMemo mocks base method.
func (m *MockMemoService) CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMemo", ctx, memo, tagNames, userID)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMemo indicates an expected call of CreateMemo.
func (mr *MockMemoServiceMockRecorder) CreateMemo(ctx, memo, tagNames, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMemo", reflect.TypeOf((*MockMemoService)(nil).CreateMemo), ctx, memo, tagNames, userID)
}

// DeleteMemo mocks base method.
func (m *MockMemoService) DeleteMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMemo", ctx, memoID, requester)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMemo indicates an expected call of DeleteMemo.
func (mr *MockMemoServiceMockRecorder) DeleteMemo(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMemo", reflect.TypeOf((*MockMemoService)(nil).DeleteMemo), ctx, memoID, requester)
}

// DeleteOrphanTags mocks base method.
func (m *MockMemoService) DeleteOrphanTags(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrphanTags", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteOrphanTags indicates an expected call of DeleteOrphanTags.
func (mr *MockMemoServiceMockRecorder) DeleteOrphanTags(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrphanTags", reflect.TypeOf((*MockMemoService)(nil).DeleteOrphanTags), arg0)
}

// GetMemo mocks base method.
func (m *MockMemoService) GetMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemo", ctx, memoID, requester)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMemo indicates an expected call of GetMemo.
func (mr *MockMemoServiceMockRecorder) GetMemo(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemo", reflect.TypeOf((*MockMemoService)(nil).GetMemo), ctx, memoID, requester)
}

// ListMemos mocks base method.
func (m *MockMemoService) ListMemos(ctx context.Context, userID uuid.UUID, tags []string, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMemos", ctx, userID, tags, sortParams, pageParams)
	ret0, _ := ret[0].([]*ent.Memo)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListMemos indicates an expected call of ListMemos.
func (mr *MockMemoServiceMockRecorder) ListMemos(ctx, userID, tags, sortParams, pageParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMemos", reflect.TypeOf((*MockMemoService)(nil).ListMemos), ctx, userID, tags, sortParams, pageParams)
}

// ListSubscribers mocks base method.
func (m *MockMemoService) ListSubscribers(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) ([]*ent.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSubscribers", ctx, memoID, requester)
	ret0, _ := ret[0].([]*ent.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSubscribers indicates an expected call of ListSubscribers.
func (mr *MockMemoServiceMockRecorder) ListSubscribers(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSubscribers", reflect.TypeOf((*MockMemoService)(nil).ListSubscribers), ctx, memoID, requester)
}

// ListTags mocks base method.
func (m *MockMemoService) ListTags(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) ([]*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTags", ctx, memoID, requester)
	ret0, _ := ret[0].([]*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTags indicates an expected call of ListTags.
func (mr *MockMemoServiceMockRecorder) ListTags(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTags", reflect.TypeOf((*MockMemoService)(nil).ListTags), ctx, memoID, requester)
}

// ReplaceTags mocks base method.
func (m *MockMemoService) ReplaceTags(ctx context.Context, memoID uuid.UUID, tagNames []string, requester *model.AppIDToken) ([]*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceTags", ctx, memoID, tagNames, requester)
	ret0, _ := ret[0].([]*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceTags indicates an expected call of ReplaceTags.
func (mr *MockMemoServiceMockRecorder) ReplaceTags(ctx, memoID, tagNames, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceTags", reflect.TypeOf((*MockMemoService)(nil).ReplaceTags), ctx, memoID, tagNames, requester)
}

// SearchTags mocks base method.
func (m *MockMemoService) SearchTags(ctx context.Context, keyword string, requester *model.AppIDToken) ([]*ent.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchTags", ctx, keyword, requester)
	ret0, _ := ret[0].([]*ent.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchTags indicates an expected call of SearchTags.
func (mr *MockMemoServiceMockRecorder) SearchTags(ctx, keyword, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTags", reflect.TypeOf((*MockMemoService)(nil).SearchTags), ctx, keyword, requester)
}

// SubscribeMemo mocks base method.
func (m *MockMemoService) SubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeMemo", ctx, memoID, requester)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeMemo indicates an expected call of SubscribeMemo.
func (mr *MockMemoServiceMockRecorder) SubscribeMemo(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeMemo", reflect.TypeOf((*MockMemoService)(nil).SubscribeMemo), ctx, memoID, requester)
}

// UnsubscribeMemo mocks base method.
func (m *MockMemoService) UnsubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsubscribeMemo", ctx, memoID, requester)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsubscribeMemo indicates an expected call of UnsubscribeMemo.
func (mr *MockMemoServiceMockRecorder) UnsubscribeMemo(ctx, memoID, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeMemo", reflect.TypeOf((*MockMemoService)(nil).UnsubscribeMemo), ctx, memoID, requester)
}

// UpdateMemo mocks base method.
func (m *MockMemoService) UpdateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, requester *model.AppIDToken) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMemo", ctx, memo, tagNames, requester)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMemo indicates an expected call of UpdateMemo.
func (mr *MockMemoServiceMockRecorder) UpdateMemo(ctx, memo, tagNames, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMemo", reflect.TypeOf((*MockMemoService)(nil).UpdateMemo), ctx, memo, tagNames, requester)
}

// UpdateMemoPublishedState mocks base method.
func (m *MockMemoService) UpdateMemoPublishedState(ctx context.Context, memoID uuid.UUID, publish bool, requester *model.AppIDToken) (*ent.Memo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMemoPublishedState", ctx, memoID, publish, requester)
	ret0, _ := ret[0].(*ent.Memo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMemoPublishedState indicates an expected call of UpdateMemoPublishedState.
func (mr *MockMemoServiceMockRecorder) UpdateMemoPublishedState(ctx, memoID, publish, requester any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMemoPublishedState", reflect.TypeOf((*MockMemoService)(nil).UpdateMemoPublishedState), ctx, memoID, publish, requester)
}
