package services

import (
	"context"
	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/platform/keycloak"

	"github.com/Nerzal/gocloak"
)

type IUserService interface {
	GetAll(ctx context.Context) ([]response.UserResponse, error)
	GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error)
	Create(ctx context.Context, userRequest request.UserRequest) error
	Delete(ctx context.Context, nickName string) error
}

type userService struct {
	keycloakService keycloak.IKeycloakService
}

func NewUserService(keycloakService keycloak.IKeycloakService) IUserService {
	return &userService{
		keycloakService: keycloakService,
	}
}

func (us *userService) GetAll(ctx context.Context) ([]response.UserResponse, error) {
	userResponse := []response.UserResponse{
		{
			Name: "andre",
		},
	}

	return userResponse, nil
}

func (us *userService) GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error) {
	var r response.UserResponse
	var err error
	return r, err
}

func (us *userService) Create(ctx context.Context, userRequest request.UserRequest) error {
	_, err := us.keycloakService.CreateUser(gocloak.User{
		Username:  userRequest.NickName,
		FirstName: userRequest.Name,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
	})

	return err
}

func (us *userService) Delete(ctx context.Context, nickName string) error {
	return nil
}
