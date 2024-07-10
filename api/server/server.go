package server

import (
	"cow_sso/middleware"

	"github.com/gin-gonic/gin"
)

type server struct {
	corsMiddleware   middleware.ICorsMiddleware
	metricMiddleWare middleware.IMetricMiddleWare
}

func New(
	middleware middleware.ICorsMiddleware,
	metricMiddleWare middleware.IMetricMiddleWare,
) *gin.Engine {
	r := gin.Default()
	r.GET("/metrics", metricMiddleWare.DefaultMetrics)

	r.Use(middleware.CorsConfig())
	r.Use(metricMiddleWare.PersonalMetrics)
	return r
}
