package server

import (
	"cow_sso/api/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	pingHandler handlers.IPingHandler
	userHandler handlers.IUserHandler
}

func NewRouter(pingHandler handlers.IPingHandler,
	userHandler handlers.IUserHandler) *Router {
	return &Router{
		pingHandler,
		userHandler,
	}
}

func (r Router) Resource(gin *gin.Engine) {
	gin.GET("/ping", r.pingHandler.Ping)

	user := gin.Group("/users")
	{
		user.GET("", r.userHandler.GetAll)
		user.GET("/:code", r.userHandler.GetByNickName)
		user.POST("", r.userHandler.Create)
		user.DELETE("/:code", r.userHandler.Delete)
	}
}
