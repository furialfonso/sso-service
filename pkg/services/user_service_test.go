package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/mocks"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	keycloakService *mocks.IKeycloakService
	restfulService  *mocks.IRestfulService
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
			name: "error CreateToken",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("", errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "error get users",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					f.keycloakService.Mock.On("GetAllUsers", mock.Anything, "ABC").Return([]*gocloak.User{}, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name: "full flow",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					id := "abcde8"
					firstName := "diego"
					lastName := "fernandez"
					email := "diego@gmail.com"
					userName := "diegof"
					f.keycloakService.Mock.On("GetAllUsers", mock.Anything, "ABC").Return([]*gocloak.User{
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
				keycloakService: &mocks.IKeycloakService{},
				restfulService:  &mocks.IRestfulService{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakService, m.restfulService)
			users, err := service.GetAll(context.Background())
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
			name:     "error CreateToken",
			nickName: "diegof",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "error GetUserByNickName",
			nickName: "diegof",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					f.keycloakService.Mock.On("GetUserByNickName", mock.Anything, "ABC", "diegof").Return([]*gocloak.User{}, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:     "full flow",
			nickName: "diegof",
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					firstName := "diego"
					lastName := "fernandez"
					email := "diego@gmail.com"
					userName := "diegof"
					f.keycloakService.Mock.On("GetUserByNickName", mock.Anything, "ABC", "diegof").Return([]*gocloak.User{
						{
							FirstName: &firstName,
							LastName:  &lastName,
							Email:     &email,
							Username:  &userName,
						},
					}, nil)
				},
			},
			outPut: response.UserResponse{
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
				keycloakService: &mocks.IKeycloakService{},
				restfulService:  &mocks.IRestfulService{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakService, m.restfulService)
			users, err := service.GetByNickName(context.Background(), tt.nickName)
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
			name: "error CreateToken",
			userRequest: request.UserRequest{
				Name:     "diego",
				LastName: "fernandez",
				Email:    "diegof@gmail.com",
				NickName: "diegof",
			},
			mocks: userMocks{
				userService: func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
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
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					role := gocloak.Role{}
					f.keycloakService.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, errors.New("some error"))
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
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					role := gocloak.Role{}
					f.keycloakService.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, nil)
					firstName := "diego"
					lastName := "fernandez"
					email := "diegof@gmail.com"
					userName := "diegof"
					f.keycloakService.Mock.On("CreateUser", mock.Anything, "ABC", &role, gocloak.User{
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
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					role := gocloak.Role{}
					f.keycloakService.Mock.On("GetRoleByID", mock.Anything, "ABC", "user").Return(&role, nil)
					firstName := "diego"
					lastName := "fernandez"
					email := "diegof@gmail.com"
					userName := "diegof"
					f.keycloakService.Mock.On("CreateUser", mock.Anything, "ABC", &role, gocloak.User{
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
				keycloakService: &mocks.IKeycloakService{},
				restfulService:  &mocks.IRestfulService{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakService, m.restfulService)
			err := service.Create(context.Background(), tt.userRequest)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		expErr error
		mocks  userMocks
		name   string
		userID string
		outPut string
	}{
		{
			name:   "error CreateToken",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:   "error GetUserByID",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					x := "diego"
					id := "1234"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakService.Mock.On("GetUserByID", mock.Anything, "ABC", "1234").Return(&user, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:   "error getting teams by user",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					x := "diego"
					id := "1234"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakService.Mock.On("GetUserByID", mock.Anything, "ABC", "1234").Return(&user, nil)
					teams := TeamsByUserResponse{}
					b, _ := json.Marshal(teams)
					f.restfulService.Mock.On("Get", mock.Anything, "http://localhost:8080/teams/users/1234", "5s").Return(b, errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:   "error user with teams",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					x := "diego"
					id := "1234"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakService.Mock.On("GetUserByID", mock.Anything, "ABC", "1234").Return(&user, nil)
					teams := TeamsByUserResponse{
						Teams: []TeamResponse{
							{
								Code: "test",
								Debt: 0,
							},
						},
					}
					b, _ := json.Marshal(teams)
					f.restfulService.Mock.On("Get", mock.Anything, "http://localhost:8080/teams/users/1234", "5s").Return(b, nil)
				},
			},
			expErr: errors.New("user 1234 has teams"),
		},
		{
			name:   "error DeleteUserByID",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					x := "diego"
					id := "1234"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakService.Mock.On("GetUserByID", mock.Anything, "ABC", "1234").Return(&user, nil)
					teams := TeamsByUserResponse{}
					b, _ := json.Marshal(teams)
					f.restfulService.Mock.On("Get", mock.Anything, "http://localhost:8080/teams/users/1234", "5s").Return(b, nil)
					f.keycloakService.Mock.On("DeleteUserByID", mock.Anything, "ABC", "1234").Return(errors.New("some error"))
				},
			},
			expErr: errors.New("some error"),
		},
		{
			name:   "full flow",
			userID: "1234",
			mocks: userMocks{
				func(f *mockUserService) {
					f.keycloakService.Mock.On("CreateToken", mock.Anything).Return("ABC", nil)
					x := "diego"
					id := "1234"
					user := gocloak.User{
						ID:       &id,
						Username: &x,
					}
					f.keycloakService.Mock.On("GetUserByID", mock.Anything, "ABC", "1234").Return(&user, nil)
					teams := TeamsByUserResponse{}
					b, _ := json.Marshal(teams)
					f.restfulService.Mock.On("Get", mock.Anything, "http://localhost:8080/teams/users/1234", "5s").Return(b, nil)
					f.keycloakService.Mock.On("DeleteUserByID", mock.Anything, "ABC", "1234").Return(nil)
				},
			},
			outPut: "diego",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{
				keycloakService: &mocks.IKeycloakService{},
				restfulService:  &mocks.IRestfulService{},
			}
			tt.mocks.userService(m)
			service := NewUserService(m.keycloakService, m.restfulService)
			users, err := service.Delete(context.Background(), tt.userID)
			if err != nil {
				assert.Equal(t, tt.expErr, err)
			}
			assert.Equal(t, tt.outPut, users)
		})
	}
}
