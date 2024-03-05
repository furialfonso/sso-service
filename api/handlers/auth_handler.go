package handlers

import (
	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	authService auth.IAuthService
}

func NewAuthHandler(authService auth.IAuthService) IAuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (a *authHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var authRequest request.AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	token, err := a.authService.Login(ctx, authRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (a *authHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	var refreshTokenRequest request.RefreshTokenRequest
	if err := c.BindJSON(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	if err := a.authService.Logout(ctx, refreshTokenRequest); err != nil {
		c.JSON(http.StatusInternalServerError, response.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "logout")
}
