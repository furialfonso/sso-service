package auth

import (
	"cow_sso/api/handlers/auth/request"
	"cow_sso/api/handlers/errors"
	"cow_sso/pkg/service/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	IsValidToken(c *gin.Context)
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
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	token, err := a.authService.Login(ctx, authRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
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
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid format",
		})
		return
	}
	if err := a.authService.Logout(ctx, refreshTokenRequest); err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "logout")
}

func (a *authHandler) IsValidToken(c *gin.Context) {
	ctx := c.Request.Context()
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "token is required",
		})
		return
	}
	token := strings.Split(auth, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, errors.ApiErrors{
			Code:    http.StatusBadRequest,
			Message: "invalid token format",
		})
		return
	}

	isValid, err := a.authService.IsValidToken(ctx, token[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ApiErrors{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if isValid {
		c.JSON(http.StatusOK, true)
	} else {
		c.JSON(http.StatusUnauthorized, false)
	}
}
