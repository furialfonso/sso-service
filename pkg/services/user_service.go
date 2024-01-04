package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/config"
	"cow_sso/pkg/platform/keycloak"
	"cow_sso/pkg/platform/restful"

	"github.com/Nerzal/gocloak/v13"
)

const (
	_roleName       = "user"
	_getTeamsByUser = "/teams/user"
)

type IUserService interface {
	GetAll(ctx context.Context) ([]response.UserResponse, error)
	GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error)
	Create(ctx context.Context, userRequest request.UserRequest) error
	Delete(ctx context.Context, userID string) (string, error)
}

type userService struct {
	keycloakService keycloak.IKeycloakService
	restfulService  restful.IRestfulService
}

func NewUserService(keycloakService keycloak.IKeycloakService,
	restfulService restful.IRestfulService,
) IUserService {
	return &userService{
		keycloakService: keycloakService,
		restfulService:  restfulService,
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
			ID:       *user.ID,
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
	user, err := us.keycloakService.GetUserByNickName(ctx, token, nickName)
	if err != nil {
		return userResponse, err
	}

	return response.UserResponse{
		ID:       *user.ID,
		Name:     *user.FirstName,
		LastName: *user.LastName,
		Email:    *user.Email,
		NickName: *user.Username,
	}, nil
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

func (us *userService) Delete(ctx context.Context, nickName string) (string, error) {
	var userName string
	token, err := us.keycloakService.CreateToken(ctx)
	if err != nil {
		return userName, err
	}
	user, err := us.keycloakService.GetUserByNickName(ctx, token, nickName)
	if err != nil {
		return userName, err
	}

	url := fmt.Sprintf("%s%s/%s", config.Get().UString("cow-api.url"), _getTeamsByUser, *user.ID)
	timeOut := config.Get().UString("cow-api.timeout")
	resp, err := us.restfulService.Get(ctx, url, timeOut)
	if err != nil {
		return userName, err
	}

	var teams TeamsByUserResponse
	err = json.Unmarshal(resp, &teams)
	if err != nil {
		return userName, err
	}

	if teams.Teams != nil {
		return userName, errors.New(fmt.Sprintf("user %s has teams", nickName))
	}

	err = us.keycloakService.DeleteUserByID(ctx, token, *user.ID)
	if err != nil {
		return userName, err
	}
	userName = *user.Username
	return userName, nil
}
