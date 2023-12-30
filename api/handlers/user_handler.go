package handlers

import (
	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	GetAll(c *gin.Context)
	GetByNickName(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	userService services.IUserService
}

func NewUserHandler(userService services.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (uh *userHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := uh.userService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: "error getting users",
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh *userHandler) GetByNickName(c *gin.Context) {
	ctx := c.Request.Context()
	nickName, exists := c.Params.Get("code")
	if !exists {
		c.JSON(http.StatusBadRequest, response.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "user's nick name is required",
		})
		return
	}
	user, err := uh.userService.GetByNickName(ctx, nickName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error getting user %s", nickName),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uh *userHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var userRequest request.UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	err := uh.userService.Create(ctx, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error creating user %s", userRequest.NickName),
		})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("user %s created", userRequest.NickName))
}

func (uh *userHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	nickName, exists := c.Params.Get("code")
	if !exists {
		c.JSON(http.StatusBadRequest, response.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "user's nick name is required",
		})
		return
	}
	err := uh.userService.Delete(ctx, nickName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting user %s", nickName),
		})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("user %s delete", nickName))
}
