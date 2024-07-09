package server

import (
	"cow_sso/api/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	pingHandler       handlers.IPingHandler
	authHandler       handlers.IAuthHandler
	userHandler       handlers.IUserHandler
	prometheusHandler handlers.IPrometheusHandler
}

func NewRouter(pingHandler handlers.IPingHandler,
	authHandler handlers.IAuthHandler,
	userHandler handlers.IUserHandler,
	prometheusHandler handlers.IPrometheusHandler,
) *Router {
	return &Router{
		pingHandler,
		authHandler,
		userHandler,
		prometheusHandler,
	}
}

func (r Router) Resource(gin *gin.Engine) {
	gin.GET("/metrics", r.prometheusHandler.GetMetrics)
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
