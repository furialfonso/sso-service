package server

import (
	"cow_sso/middleware"

	"github.com/gin-gonic/gin"
)

type server struct {
	middleware middleware.ICorsConfig
}

func New(middleware middleware.ICorsConfig) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsConfig())
	return r
}
