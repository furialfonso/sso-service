package services

import (
	"context"
	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
)

type IUserService interface {
	GetAll(ctx context.Context) ([]response.UserResponse, error)
	GetByNickName(ctx context.Context, nickName string) (response.UserResponse, error)
	Create(ctx context.Context, userRequest request.UserRequest) error
	Delete(ctx context.Context, nickName string) error
}

type userService struct{}

func NewUserService() IUserService {
	return &userService{}
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
	return nil
}

func (us *userService) Delete(ctx context.Context, nickName string) error {
	return nil
}
