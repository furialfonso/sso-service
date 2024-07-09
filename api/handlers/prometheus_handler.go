package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type IPrometheusHandler interface {
	GetMetrics(c *gin.Context)
}

type prometheusHandler struct {
}

func NewPrometheusHandler() IPrometheusHandler {
	return &prometheusHandler{}
}

func (p *prometheusHandler) GetMetrics(c *gin.Context) {
	h := promhttp.Handler()
	h.ServeHTTP(c.Writer, c.Request)
}
