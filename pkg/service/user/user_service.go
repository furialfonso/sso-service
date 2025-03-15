package user

import (
	"context"
	"errors"
	"fmt"

	"cow_sso/api/handlers/user/request"
	"cow_sso/api/handlers/user/response"
	"cow_sso/pkg/integration/keycloak"
	"cow_sso/pkg/integration/team"

	"github.com/Nerzal/gocloak/v13"
)

const (
	_roleName       = "user"
	_getTeamsByUser = "/teams/user"
)

type IUserService interface {
	GetAll(ctx context.Context, token string) ([]response.UserResponse, error)
	GetByNickName(ctx context.Context, token string, nickName string) (response.UserResponse, error)
	Create(ctx context.Context, token string, userRequest request.UserRequest) error
	Delete(ctx context.Context, token string, userID string) (string, error)
}

type userService struct {
	keycloakClient keycloak.IKeycloakClient
	teamClient     team.ITeamClient
}

func NewUserService(keycloakClient keycloak.IKeycloakClient,
	teamClient team.ITeamClient,
) IUserService {
	return &userService{
		keycloakClient: keycloakClient,
		teamClient:     teamClient,
	}
}

func (us *userService) GetAll(ctx context.Context, token string) ([]response.UserResponse, error) {
	var userResponse []response.UserResponse
	users, err := us.keycloakClient.GetAllUsers(ctx, token)
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

func (us *userService) GetByNickName(ctx context.Context, token string, nickName string) (response.UserResponse, error) {
	var userResponse response.UserResponse
	user, err := us.keycloakClient.GetUserByNickName(ctx, token, nickName)
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

func (us *userService) Create(ctx context.Context, token string, userRequest request.UserRequest) error {
	role, err := us.keycloakClient.GetRoleByID(ctx, token, _roleName)
	if err != nil {
		return err
	}
	_, err = us.keycloakClient.CreateUser(ctx, token, role, gocloak.User{
		Username:  &userRequest.NickName,
		FirstName: &userRequest.Name,
		LastName:  &userRequest.LastName,
		Email:     &userRequest.Email,
	})
	return err
}

func (us *userService) Delete(ctx context.Context, token string, nickName string) (string, error) {
	var userName string
	user, err := us.keycloakClient.GetUserByNickName(ctx, token, nickName)
	if err != nil {
		return userName, err
	}

	teams, err := us.teamClient.GetTeamsByUser(ctx, *user.ID)
	if err != nil {
		return userName, err
	}

	if teams.Teams != nil {
		return userName, errors.New(fmt.Sprintf("user %s has teams", nickName))
	}

	err = us.keycloakClient.DeleteUserByID(ctx, token, *user.ID)
	if err != nil {
		return userName, err
	}
	userName = *user.Username
	return userName, nil
}
