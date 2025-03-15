package server

import (
	"cow_sso/api/handlers"
	"cow_sso/api/handlers/auth"
	"cow_sso/api/handlers/user"

	"github.com/gin-gonic/gin"
)

type Router struct {
	pingHandler handlers.IPingHandler
	authHandler auth.IAuthHandler
	userHandler user.IUserHandler
}

func NewRouter(pingHandler handlers.IPingHandler,
	authHandler auth.IAuthHandler,
	userHandler user.IUserHandler,

) *Router {
	return &Router{
		pingHandler,
		authHandler,
		userHandler,
	}
}

func (r Router) Resource(gin *gin.Engine) {
	gin.GET("/ping", r.pingHandler.Ping)
	auth := gin.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/logout", r.authHandler.Logout)
		auth.POST("/valid-token", r.authHandler.IsValidToken)
	}
	user := gin.Group("/users")
	{
		user.GET("", r.userHandler.GetAll)
		user.GET("/:code", r.userHandler.GetByNickName)
		user.POST("", r.userHandler.Create)
		user.DELETE("/:code", r.userHandler.Delete)
	}
}
