package user

import (
	"fmt"
	"net/http"

	"cow_sso/api/handlers/errors"
	"cow_sso/api/handlers/user/request"
	"cow_sso/pkg/service/user"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	GetAll(c *gin.Context)
	GetByNickName(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	userService user.IUserService
}

func NewUserHandler(userService user.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (uh *userHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, errors.ApiErrors{
			Code:    http.StatusUnauthorized,
			Message: "token is required",
		})
		return
	}

	users, err := uh.userService.GetAll(ctx, token[7:])
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: "error getting users",
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh *userHandler) GetByNickName(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, errors.ApiErrors{
			Code:    http.StatusUnauthorized,
			Message: "token is required",
		})
		return
	}

	nickName, exists := c.Params.Get("code")
	if !exists {
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "user's nick name is required",
		})
		return
	}

	user, err := uh.userService.GetByNickName(ctx, token[7:], nickName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error getting user %s, err: %s", nickName, err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uh *userHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, errors.ApiErrors{
			Code:    http.StatusUnauthorized,
			Message: "token is required",
		})
		return
	}

	var userRequest request.UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	err := uh.userService.Create(ctx, token[7:], userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error creating user %s, err: %s", userRequest.NickName, err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("user %s created", userRequest.NickName))
}

func (uh *userHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, errors.ApiErrors{
			Code:    http.StatusUnauthorized,
			Message: "token is required",
		})
		return
	}

	nickName, exists := c.Params.Get("code")
	if !exists {
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "user's nick name is required",
		})
		return
	}

	userName, err := uh.userService.Delete(ctx, token[7:], nickName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting user: %s, err: %s", nickName, err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("user %s delete", userName))
}
