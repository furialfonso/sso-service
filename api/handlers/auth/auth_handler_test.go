package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cow_sso/api/handlers/auth/request"
	"cow_sso/api/handlers/auth/response"
	"cow_sso/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthHandler struct {
	authService *mocks.IAuthService
}

type authMocks struct {
	authHandler func(f *mockAuthHandler)
}

func Test_Login(t *testing.T) {
	tests := []struct {
		mocks       authMocks
		name        string
		authRequest any
		expCode     int
	}{
		{
			name:        "error in bind json",
			authRequest: "invalid format",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "error in login",
			authRequest: request.AuthRequest{
				User:     "user",
				Password: "password",
			},
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("Login", mock.Anything, request.AuthRequest{
						User:     "user",
						Password: "password",
					}).Return(response.AuthResponse{}, errors.New("some error"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			authRequest: request.AuthRequest{
				User:     "user",
				Password: "password",
			},
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("Login", mock.Anything, request.AuthRequest{
						User:     "user",
						Password: "password",
					}).Return(response.AuthResponse{
						Token:            "token",
						ExpiresIn:        3600,
						RefreshToken:     "refresh_token",
						RefreshExpiresIn: 3600,
					}, nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockAuthHandler{
				&mocks.IAuthService{},
			}
			tc.mocks.authHandler(ms)
			handler := NewAuthHandler(ms.authService)
			url := "/auth/login"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.POST(url, func(ctx *gin.Context) {
				handler.Login(ctx)
			})
			res := httptest.NewRecorder()
			b, _ := json.Marshal(tc.authRequest)
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewBuffer(b)))
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}

func Test_Logout(t *testing.T) {
	tests := []struct {
		mocks               authMocks
		name                string
		refreshTokenRequest any
		expCode             int
	}{
		{
			name:                "error in bind json",
			refreshTokenRequest: "invalid format",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "error in logout",
			refreshTokenRequest: request.RefreshTokenRequest{
				RefreshToken: "refresh_token",
			},
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("Logout", mock.Anything, request.RefreshTokenRequest{
						RefreshToken: "refresh_token",
					}).Return(errors.New("some error"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			refreshTokenRequest: request.RefreshTokenRequest{
				RefreshToken: "refresh_token",
			},
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("Logout", mock.Anything, request.RefreshTokenRequest{
						RefreshToken: "refresh_token",
					}).Return(nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockAuthHandler{
				&mocks.IAuthService{},
			}
			tc.mocks.authHandler(ms)
			handler := NewAuthHandler(ms.authService)
			url := "/auth/logout"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.POST(url, func(ctx *gin.Context) {
				handler.Logout(ctx)
			})
			res := httptest.NewRecorder()
			b, _ := json.Marshal(tc.refreshTokenRequest)
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewBuffer(b)))
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}

func Test_IsValidToken(t *testing.T) {
	tests := []struct {
		mocks       authMocks
		name        string
		accessToken string
		expCode     int
	}{
		{
			name: "token is required",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:        "error token format",
			accessToken: "Bearertoken",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:        "error toekn",
			accessToken: "Bearer token",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("IsValidToken", mock.Anything, "token").Return(false, errors.New("some error"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:        "token invalid",
			accessToken: "Bearer token",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("IsValidToken", mock.Anything, "token").Return(false, nil)
				},
			},
			expCode: http.StatusUnauthorized,
		},
		{
			name:        "success",
			accessToken: "Bearer token",
			mocks: authMocks{
				authHandler: func(f *mockAuthHandler) {
					f.authService.On("IsValidToken", mock.Anything, "token").Return(true, nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockAuthHandler{
				&mocks.IAuthService{},
			}
			tc.mocks.authHandler(ms)
			handler := NewAuthHandler(ms.authService)
			url := "/auth/valid-token"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.POST(url, func(ctx *gin.Context) {
				handler.IsValidToken(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, url, nil)
			if tc.accessToken != "" {
				req.Header.Set("Authorization", tc.accessToken)
			}
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}
