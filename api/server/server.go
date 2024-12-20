package server

import (
	"cow_sso/middleware"

	"github.com/gin-gonic/gin"
)

type server struct {
	metricMiddleWare middleware.IMetricMiddleWare
}

func New(
	metricMiddleWare middleware.IMetricMiddleWare,
) *gin.Engine {
	r := gin.Default()
	r.GET("/metrics", metricMiddleWare.DefaultMetrics)

	r.Use(metricMiddleWare.PersonalMetrics)
	return r
}
