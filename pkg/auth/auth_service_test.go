package auth

import (
	"context"
	"errors"
	"testing"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/mocks"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	keycloakService *mocks.IKeycloakService
}

type authMocks struct {
	authService func(f *mockAuthService)
}

func Test_Login(t *testing.T) {
	tests := []struct {
		name        string
		authRequest request.AuthRequest
		mocks       authMocks
		outPut      response.AuthResponse
		expErr      error
	}{
		{
			name: "Login fail",
			authRequest: request.AuthRequest{
				User:     "user",
				Password: "password",
			},
			mocks: authMocks{
				func(f *mockAuthService) {
					x := gocloak.JWT{}
					f.keycloakService.On("Login", mock.Anything, "user", "password").Return(&x, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "Login success",
			authRequest: request.AuthRequest{
				User:     "user",
				Password: "password",
			},
			mocks: authMocks{
				func(f *mockAuthService) {
					x := gocloak.JWT{
						AccessToken:      "token",
						ExpiresIn:        3600,
						RefreshToken:     "refresh_token",
						RefreshExpiresIn: 3600,
					}
					f.keycloakService.On("Login", mock.Anything, "user", "password").Return(&x, nil)
				},
			},
			outPut: response.AuthResponse{
				Token:            "token",
				ExpiresIn:        3600,
				RefreshToken:     "refresh_token",
				RefreshExpiresIn: 3600,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockAuthService{
				keycloakService: &mocks.IKeycloakService{},
			}
			tt.mocks.authService(m)
			service := NewAuthService(m.keycloakService)
			auth, err := service.Login(context.Background(), tt.authRequest)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, auth)
		})
	}
}

func Test_Logout(t *testing.T) {
	tests := []struct {
		name                string
		refreshTokenRequest request.RefreshTokenRequest
		mocks               authMocks
		expErr              error
	}{
		{
			name: "Logout fail",
			refreshTokenRequest: request.RefreshTokenRequest{
				RefreshToken: "refresh_token",
			},
			mocks: authMocks{
				func(f *mockAuthService) {
					f.keycloakService.On("Logout", mock.Anything, "refresh_token").Return(errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "Logout success",
			refreshTokenRequest: request.RefreshTokenRequest{
				RefreshToken: "refresh_token",
			},
			mocks: authMocks{
				func(f *mockAuthService) {
					f.keycloakService.On("Logout", mock.Anything, "refresh_token").Return(nil)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockAuthService{
				keycloakService: &mocks.IKeycloakService{},
			}
			tt.mocks.authService(m)
			service := NewAuthService(m.keycloakService)
			err := service.Logout(context.Background(), tt.refreshTokenRequest)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func Test_IsValidToken(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		mocks       authMocks
		expErr      error
		outPut      bool
	}{
		{
			name:        "error valid token",
			accessToken: "abc",
			mocks: authMocks{
				func(f *mockAuthService) {
					f.keycloakService.On("IsValidToken", mock.Anything, "abc").Return(false, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:        "valid token",
			accessToken: "abc",
			mocks: authMocks{
				func(f *mockAuthService) {
					f.keycloakService.On("IsValidToken", mock.Anything, "abc").Return(true, nil)
				},
			},
			outPut: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockAuthService{
				keycloakService: &mocks.IKeycloakService{},
			}
			tt.mocks.authService(m)
			service := NewAuthService(m.keycloakService)
			resp, err := service.IsValidToken(context.Background(), tt.accessToken)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, resp)
		})
	}
}
