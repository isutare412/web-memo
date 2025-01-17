// Code generated by mockery v2.46.1. DO NOT EDIT.

package mockport

import (
	context "context"

	model "github.com/isutare412/web-memo/api/internal/core/model"
	mock "github.com/stretchr/testify/mock"
)

// MockGoogleClient is an autogenerated mock type for the GoogleClient type
type MockGoogleClient struct {
	mock.Mock
}

type MockGoogleClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGoogleClient) EXPECT() *MockGoogleClient_Expecter {
	return &MockGoogleClient_Expecter{mock: &_m.Mock}
}

// ExchangeAuthCode provides a mock function with given fields: ctx, code, redirectURI
func (_m *MockGoogleClient) ExchangeAuthCode(ctx context.Context, code string, redirectURI string) (model.GoogleTokenResponse, error) {
	ret := _m.Called(ctx, code, redirectURI)

	if len(ret) == 0 {
		panic("no return value specified for ExchangeAuthCode")
	}

	var r0 model.GoogleTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (model.GoogleTokenResponse, error)); ok {
		return rf(ctx, code, redirectURI)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) model.GoogleTokenResponse); ok {
		r0 = rf(ctx, code, redirectURI)
	} else {
		r0 = ret.Get(0).(model.GoogleTokenResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, code, redirectURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGoogleClient_ExchangeAuthCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExchangeAuthCode'
type MockGoogleClient_ExchangeAuthCode_Call struct {
	*mock.Call
}

// ExchangeAuthCode is a helper method to define mock.On call
//   - ctx context.Context
//   - code string
//   - redirectURI string
func (_e *MockGoogleClient_Expecter) ExchangeAuthCode(ctx interface{}, code interface{}, redirectURI interface{}) *MockGoogleClient_ExchangeAuthCode_Call {
	return &MockGoogleClient_ExchangeAuthCode_Call{Call: _e.mock.On("ExchangeAuthCode", ctx, code, redirectURI)}
}

func (_c *MockGoogleClient_ExchangeAuthCode_Call) Run(run func(ctx context.Context, code string, redirectURI string)) *MockGoogleClient_ExchangeAuthCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockGoogleClient_ExchangeAuthCode_Call) Return(_a0 model.GoogleTokenResponse, _a1 error) *MockGoogleClient_ExchangeAuthCode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGoogleClient_ExchangeAuthCode_Call) RunAndReturn(run func(context.Context, string, string) (model.GoogleTokenResponse, error)) *MockGoogleClient_ExchangeAuthCode_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGoogleClient creates a new instance of MockGoogleClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGoogleClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGoogleClient {
	mock := &MockGoogleClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
