// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	request "cow_sso/api/handlers/auth/request"

	mock "github.com/stretchr/testify/mock"

	response "cow_sso/api/handlers/auth/response"
)

// IAuthService is an autogenerated mock type for the IAuthService type
type IAuthService struct {
	mock.Mock
}

// IsValidToken provides a mock function with given fields: ctx, accessToken
func (_m *IAuthService) IsValidToken(ctx context.Context, accessToken string) (bool, error) {
	ret := _m.Called(ctx, accessToken)

	if len(ret) == 0 {
		panic("no return value specified for IsValidToken")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, accessToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, accessToken)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, authRequest
func (_m *IAuthService) Login(ctx context.Context, authRequest request.AuthRequest) (response.AuthResponse, error) {
	ret := _m.Called(ctx, authRequest)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 response.AuthResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.AuthRequest) (response.AuthResponse, error)); ok {
		return rf(ctx, authRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.AuthRequest) response.AuthResponse); ok {
		r0 = rf(ctx, authRequest)
	} else {
		r0 = ret.Get(0).(response.AuthResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.AuthRequest) error); ok {
		r1 = rf(ctx, authRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: ctx, refreshTokenRequest
func (_m *IAuthService) Logout(ctx context.Context, refreshTokenRequest request.RefreshTokenRequest) error {
	ret := _m.Called(ctx, refreshTokenRequest)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, request.RefreshTokenRequest) error); ok {
		r0 = rf(ctx, refreshTokenRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIAuthService creates a new instance of IAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthService {
	mock := &IAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
