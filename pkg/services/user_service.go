package services

import (
	"context"
	"errors"
	"fmt"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/platform/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

const (
	_roleName = "user"
)

type IUserService interface {
	GetAll(ctx context.Context) ([]response.UserResponse, error)
	GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error)
	Create(ctx context.Context, userRequest request.UserRequest) error
	Delete(ctx context.Context, userID string) (string, error)
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
	var userResponse []response.UserResponse
	token, err := us.keycloakService.CreateToken(ctx)
	if err != nil {
		return userResponse, err
	}
	users, err := us.keycloakService.GetAllUsers(ctx, token)
	if err != nil {
		return userResponse, err
	}
	for _, user := range users {
		userResponse = append(userResponse, response.UserResponse{
			Name:     *user.FirstName,
			LastName: *user.LastName,
			Email:    *user.Email,
			NickName: *user.Username,
		})
	}
	return userResponse, nil
}

func (us *userService) GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error) {
	var userResponse response.UserResponse
	token, err := us.keycloakService.CreateToken(ctx)
	if err != nil {
		return userResponse, err
	}
	users, err := us.keycloakService.GetUserByNickName(ctx, token, nickName)
	if err != nil {
		return userResponse, err
	}
	if len(users) == 0 {
		return userResponse, errors.New(fmt.Sprintf("user %s doesn't exist", nickName))
	}
	for _, user := range users {
		userResponse = response.UserResponse{
			Name:     *user.FirstName,
			LastName: *user.LastName,
			Email:    *user.Email,
			NickName: *user.Username,
		}
		break
	}
	return userResponse, nil
}

func (us *userService) Create(ctx context.Context, userRequest request.UserRequest) error {
	token, err := us.keycloakService.CreateToken(ctx)
	if err != nil {
		return err
	}
	role, err := us.keycloakService.GetRoleByID(ctx, token, _roleName)
	if err != nil {
		return err
	}
	_, err = us.keycloakService.CreateUser(ctx, token, role, gocloak.User{
		Username:  &userRequest.NickName,
		FirstName: &userRequest.Name,
		LastName:  &userRequest.LastName,
		Email:     &userRequest.Email,
	})

	return err
}

func (us *userService) Delete(ctx context.Context, userID string) (string, error) {
	var userName string
	token, err := us.keycloakService.CreateToken(ctx)
	if err != nil {
		return userName, err
	}
	user, err := us.keycloakService.GetUserByID(ctx, token, userID)
	if err != nil {
		return userName, err
	}

	err = us.keycloakService.DeleteUserByID(ctx, token, userID)
	if err != nil {
		return userName, err
	}
	userName = *user.Username
	return userName, nil
}
