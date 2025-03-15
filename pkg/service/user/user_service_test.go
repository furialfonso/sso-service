package user

import (
	"context"
	"errors"
	"testing"

	"cow_sso/api/handlers/user/request"
	"cow_sso/api/handlers/user/response"
	"cow_sso/mocks"
	"cow_sso/pkg/repository/team/dto"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	keycloakRepository *mocks.IKeycloakRepository
	teamRepository     *mocks.ITeamRepository
}

type userMocks struct {
	userService func(f *mockUserService)
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		expErr error
		mocks  userMocks
		name   string
		outPut []response.UserResponse
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakRepository.Mock.On("GetAllUsers", mock.Anything, "ABC").Return([]*gocloak.User{}, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "full flow",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					id := "abcde8"
					firstName := "diego"
					lastName := "fernandez"
					email := "diego@gmail.com"
					userName := "diegof"
					f.keycloakRepository.Mock.On("GetAllUsers", mock.Anything, "ABC").Return([]*gocloak.User{
						{
							ID:        &id,
							FirstName: &firstName,
							LastName:  &lastName,
							Email:     &email,
							Username:  &userName,
						},
					}, nil)
				},
			},
			outPut: []response.UserResponse{
				{
					ID:       "abcde8",
					Name:     "diego",
					LastName: "fernandez",
					Email:    "diego@gmail.com",
					NickName: "diegof",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{
				keycloakRepository: &mocks.IKeycloakRepository{},
				teamRepository:     &mocks.ITeamRepository{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakRepository, m.teamRepository)
			users, err := service.GetAll(context.Background(), mock.Anything)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, users)
		})
	}
}

func Test_GetByNickName(t *testing.T) {
	tests := []struct {
		expErr   error
		mocks    userMocks
		outPut   response.UserResponse
		name     string
		nickName string
	}{
		{
			name:     "error GetUserByNickName",
			nickName: "diegof",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					user := gocloak.User{}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "diegof").Return(&user, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "full flow",
			nickName: "diegof",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					id := "abcde8"
					firstName := "diego"
					lastName := "fernandez"
					email := "diego@gmail.com"
					userName := "diegof"
					user := gocloak.User{
						ID:        &id,
						FirstName: &firstName,
						LastName:  &lastName,
						Email:     &email,
						Username:  &userName,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "diegof").Return(&user, nil)
				},
			},
			outPut: response.UserResponse{
				ID:       "abcde8",
				Name:     "diego",
				LastName: "fernandez",
				Email:    "diego@gmail.com",
				NickName: "diegof",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{
				keycloakRepository: &mocks.IKeycloakRepository{},
				teamRepository:     &mocks.ITeamRepository{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakRepository, m.teamRepository)
			users, err := service.GetByNickName(context.Background(), mock.Anything, tt.nickName)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, users)
		})
	}
}

func Test_Create(t *testing.T) {
	tests := []struct {
		expErr      error
		mocks       userMocks
		userRequest request.UserRequest
		name        string
	}{
		{
			name: "error GetRoleByID",
			userRequest: request.UserRequest{
				Name:     "diego",
				LastName: "fernandez",
				Email:    "diegof@gmail.com",
				NickName: "diegof",
			},
			mocks: userMocks{
				userService: func(f *mockUserService) {
					role := gocloak.Role{}
					f.keycloakRepository.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "error CreateUser",
			userRequest: request.UserRequest{
				Name:     "diego",
				LastName: "fernandez",
				Email:    "diegof@gmail.com",
				NickName: "diegof",
			},
			mocks: userMocks{
				userService: func(f *mockUserService) {
					role := gocloak.Role{}
					f.keycloakRepository.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, nil)
					firstName := "diego"
					lastName := "fernandez"
					email := "diegof@gmail.com"
					userName := "diegof"
					f.keycloakRepository.Mock.On("CreateUser", mock.Anything, "ABC", &role, gocloak.User{
						FirstName: &firstName,
						LastName:  &lastName,
						Email:     &email,
						Username:  &userName,
					}).Return("", errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "full flow",
			userRequest: request.UserRequest{
				Name:     "diego",
				LastName: "fernandez",
				Email:    "diegof@gmail.com",
				NickName: "diegof",
			},
			mocks: userMocks{
				userService: func(f *mockUserService) {
					role := gocloak.Role{}
					f.keycloakRepository.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, nil)
					firstName := "diego"
					lastName := "fernandez"
					email := "diegof@gmail.com"
					userName := "diegof"
					f.keycloakRepository.Mock.On("CreateUser", mock.Anything, "ABC", &role, gocloak.User{
						FirstName: &firstName,
						LastName:  &lastName,
						Email:     &email,
						Username:  &userName,
					}).Return("123", nil)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{
				keycloakRepository: &mocks.IKeycloakRepository{},
				teamRepository:     &mocks.ITeamRepository{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakRepository, m.teamRepository)
			err := service.Create(context.Background(), mock.Anything, tt.userRequest)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		expErr   error
		mocks    userMocks
		name     string
		token    string
		nickName string
		outPut   string
	}{
		{
			name:     "error GetUserByID",
			nickName: "1234",
			token:    "ABC",
			mocks: userMocks{
				func(f *mockUserService) {
					x := "diego"
					id := "AXYZT"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "1234").Return(&user, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "error user with teams unmarshal",
			nickName: "1234",
			token:    "ABC",
			mocks: userMocks{
				func(f *mockUserService) {
					x := "diego"
					id := "AXYZT"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "1234").Return(&user, nil)
					f.teamRepository.Mock.On("GetTeamsByUser", mock.Anything, "AXYZT").Return(dto.TeamsByUserResponse{}, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "error user with teams",
			token:    "ABC",
			nickName: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					x := "diego"
					id := "AXYZT"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "1234").Return(&user, nil)
					f.teamRepository.Mock.On("GetTeamsByUser", mock.Anything, "AXYZT").Return(dto.TeamsByUserResponse{
						Teams: []dto.TeamResponse{
							{
								Code: "xs",
								Debt: 0,
							},
						},
					}, nil)
				},
			},
			expErr: errors.New("user 1234 has teams"),
		},
		{
			name:     "error DeleteUserByID",
			token:    "ABC",
			nickName: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					x := "diego"
					id := "AXYZT"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "1234").Return(&user, nil)
					f.teamRepository.Mock.On("GetTeamsByUser", mock.Anything, "AXYZT").Return(dto.TeamsByUserResponse{}, nil)
					f.keycloakRepository.Mock.On("DeleteUserByID", mock.Anything, "ABC", "AXYZT").Return(errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "full flow",
			token:    "ABC",
			nickName: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					x := "diego"
					id := "AXYZT"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakRepository.Mock.On("GetUserByNickName", mock.Anything, "ABC", "1234").Return(&user, nil)
					f.teamRepository.Mock.On("GetTeamsByUser", mock.Anything, "AXYZT").Return(dto.TeamsByUserResponse{}, nil)
					f.keycloakRepository.Mock.On("DeleteUserByID", mock.Anything, "ABC", "AXYZT").Return(nil)
				},
			},
			outPut: "diego",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{
				keycloakRepository: &mocks.IKeycloakRepository{},
				teamRepository:     &mocks.ITeamRepository{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakRepository, m.teamRepository)
			users, err := service.Delete(context.Background(), tt.token, tt.nickName)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, users)
		})
	}
}
