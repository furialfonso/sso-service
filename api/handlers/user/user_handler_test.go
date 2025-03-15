package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cow_sso/api/handlers/user/request"
	"cow_sso/api/handlers/user/response"
	"cow_sso/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserHandler struct {
	userService *mocks.IUserService
}

type userMocks struct {
	userHandler func(f *mockUserHandler)
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		mocks       userMocks
		name        string
		token       string
		expNickName int
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expNickName: http.StatusUnauthorized,
		},
		{
			name:  "error get users",
			token: "Bearer ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetAll", mock.Anything, "ABC").Return([]response.UserResponse{}, errors.New("error x"))
				},
			},
			expNickName: http.StatusInternalServerError,
		},
		{
			name:  "full flow",
			token: "Bearer ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetAll", mock.Anything, "ABC").Return([]response.UserResponse{}, nil)
				},
			},
			expNickName: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.GET(url, func(ctx *gin.Context) {
				handler.GetAll(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expNickName, res.Code)
		})
	}
}

func Test_GetByNickName(t *testing.T) {
	tests := []struct {
		mocks    userMocks
		name     string
		nickName string
		token    string
		expCode  int
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusUnauthorized,
		},
		{
			name:     "code not sending",
			token:    "Bearer ABC",
			nickName: "",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:     "error getting user by id",
			token:    "Bearer ABC",
			nickName: "ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetByNickName", mock.Anything, "ABC", "ABC").Return(response.UserResponse{}, errors.New("x"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:     "full flow",
			token:    "Bearer ABC",
			nickName: "ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetByNickName", mock.Anything, "ABC", "ABC").Return(response.UserResponse{
						NickName: "ABC",
					}, nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.GET(url, func(ctx *gin.Context) {
				if tc.nickName != "" {
					ctx.AddParam("code", tc.nickName)
				}
				handler.GetByNickName(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}

func Test_Create(t *testing.T) {
	tests := []struct {
		input   interface{}
		mocks   userMocks
		name    string
		token   string
		expCode int
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusUnauthorized,
		},
		{
			name:  "error on input",
			input: "ABC",
			token: "Bearer ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:  "error creating user",
			token: "Bearer ABC",
			input: request.UserRequest{
				Name:     "a",
				NickName: "c",
				Email:    "d",
				LastName: "e",
			},
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Create", mock.Anything, "ABC", request.UserRequest{
						Name:     "a",
						NickName: "c",
						Email:    "d",
						LastName: "e",
					}).Return(errors.New("error creating user"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:  "full flow",
			token: "Bearer ABC",
			input: request.UserRequest{
				Name:     "a",
				NickName: "c",
				Email:    "d",
				LastName: "e",
			},
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Create", mock.Anything, "ABC", request.UserRequest{
						Name:     "a",
						NickName: "c",
						Email:    "d",
						LastName: "e",
					}).Return(nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users/"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.POST(url, func(ctx *gin.Context) {
				handler.Create(ctx)
			})
			res := httptest.NewRecorder()
			b, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewBuffer(b)))
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		mocks   userMocks
		name    string
		userID  string
		token   string
		expCode int
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusUnauthorized,
		},
		{
			name:  "nick name isnt present",
			token: "Bearer ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:   "nick name not found",
			token:  "Bearer ABC",
			userID: "abc",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Delete", mock.Anything, "ABC", "abc").Return("", errors.New("nick name not found"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "full flow",
			token:  "Bearer ABC",
			userID: "abc",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Delete", mock.Anything, "ABC", "abc").Return("test", nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			url := "/users"
			engine.DELETE(url, func(ctx *gin.Context) {
				if tc.userID != "" {
					ctx.AddParam("code", tc.userID)
				}
				handler.Delete(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}
