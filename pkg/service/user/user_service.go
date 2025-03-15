package user

import (
	"context"
	"errors"
	"fmt"

	"cow_sso/api/handlers/user/request"
	"cow_sso/api/handlers/user/response"
	"cow_sso/pkg/repository/keycloak"
	"cow_sso/pkg/repository/team"

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
	keycloakRepository keycloak.IKeycloakRepository
	teamRepository     team.ITeamRepository
}

func NewUserService(keycloakRepository keycloak.IKeycloakRepository,
	teamRepository team.ITeamRepository,
) IUserService {
	return &userService{
		keycloakRepository: keycloakRepository,
		teamRepository:     teamRepository,
	}
}

func (us *userService) GetAll(ctx context.Context, token string) ([]response.UserResponse, error) {
	var userResponse []response.UserResponse
	users, err := us.keycloakRepository.GetAllUsers(ctx, token)
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
	user, err := us.keycloakRepository.GetUserByNickName(ctx, token, nickName)
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
	role, err := us.keycloakRepository.GetRoleByID(ctx, token, _roleName)
	if err != nil {
		return err
	}
	_, err = us.keycloakRepository.CreateUser(ctx, token, role, gocloak.User{
		Username:  &userRequest.NickName,
		FirstName: &userRequest.Name,
		LastName:  &userRequest.LastName,
		Email:     &userRequest.Email,
	})
	return err
}

func (us *userService) Delete(ctx context.Context, token string, nickName string) (string, error) {
	var userName string
	user, err := us.keycloakRepository.GetUserByNickName(ctx, token, nickName)
	if err != nil {
		return userName, err
	}

	teams, err := us.teamRepository.GetTeamsByUser(ctx, *user.ID)
	if err != nil {
		return userName, err
	}

	if teams.Teams != nil {
		return userName, errors.New(fmt.Sprintf("user %s has teams", nickName))
	}

	err = us.keycloakRepository.DeleteUserByID(ctx, token, *user.ID)
	if err != nil {
		return userName, err
	}
	userName = *user.Username
	return userName, nil
}
