package auth

import (
	"context"

	"cow_sso/api/handlers/auth/request"
	"cow_sso/api/handlers/auth/response"
	"cow_sso/pkg/integration/keycloak"
)

type IAuthService interface {
	Login(ctx context.Context, authRequest request.AuthRequest) (response.AuthResponse, error)
	Logout(ctx context.Context, refreshTokenRequest request.RefreshTokenRequest) error
	IsValidToken(ctx context.Context, accessToken string) (bool, error)
}

type authService struct {
	keycloakClient keycloak.IKeycloakClient
}

func NewAuthService(keycloakClient keycloak.IKeycloakClient) IAuthService {
	return &authService{
		keycloakClient: keycloakClient,
	}
}

func (a *authService) Login(ctx context.Context, authRequest request.AuthRequest) (response.AuthResponse, error) {
	var authResponse response.AuthResponse
	token, err := a.keycloakClient.Login(ctx, authRequest.User, authRequest.Password)
	if err != nil {
		return authResponse, err
	}
	authResponse.Token = token.AccessToken
	authResponse.ExpiresIn = token.ExpiresIn
	authResponse.RefreshToken = token.RefreshToken
	authResponse.RefreshExpiresIn = token.RefreshExpiresIn
	return authResponse, nil
}

func (a *authService) Logout(ctx context.Context, refreshTokenRequest request.RefreshTokenRequest) error {
	return a.keycloakClient.Logout(ctx, refreshTokenRequest.RefreshToken)
}

func (a *authService) IsValidToken(ctx context.Context, accessToken string) (bool, error) {
	return a.keycloakClient.IsValidToken(ctx, accessToken)
}
