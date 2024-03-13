package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ICorsConfig interface {
	CorsConfig() gin.HandlerFunc
}

type corsConfig struct{}

func NewCorsConfig() ICorsConfig {
	return &corsConfig{}
}

func (cc *corsConfig) CorsConfig() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	return cors.New(config)
}
