package auth

import (
	"context"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/platform/keycloak"
)

type IAuthService interface {
	Login(ctx context.Context, authRequest request.AuthRequest) (response.AuthResponse, error)
	Logout(ctx context.Context, refreshTokenRequest request.RefreshTokenRequest) error
	IsValidToken(ctx context.Context, accessToken string) (bool, error)
}

type authService struct {
	keycloakService keycloak.IKeycloakService
}

func NewAuthService(keycloakService keycloak.IKeycloakService) IAuthService {
	return &authService{
		keycloakService: keycloakService,
	}
}

func (a *authService) Login(ctx context.Context, authRequest request.AuthRequest) (response.AuthResponse, error) {
	var authResponse response.AuthResponse
	token, err := a.keycloakService.Login(ctx, authRequest.User, authRequest.Password)
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
	return a.keycloakService.Logout(ctx, refreshTokenRequest.RefreshToken)
}

func (a *authService) IsValidToken(ctx context.Context, accessToken string) (bool, error) {
	return a.keycloakService.IsValidToken(ctx, accessToken)
}
