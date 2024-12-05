package middleware

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	instance = os.Getenv("INSTANCE_ID")

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status", "instance"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "instance"},
	)
)

type IMetricMiddleWare interface {
	DefaultMetrics(c *gin.Context)
	PersonalMetrics(c *gin.Context)
}

type metricMiddleWare struct{}

func NewMetricMiddleWare() IMetricMiddleWare {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
	return &metricMiddleWare{}
}

func (p *metricMiddleWare) DefaultMetrics(c *gin.Context) {
	h := promhttp.Handler()
	h.ServeHTTP(c.Writer, c.Request)
}

func (m *metricMiddleWare) PersonalMetrics(c *gin.Context) {
	path := c.FullPath()
	method := c.Request.Method
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(method, path, "", instance))
	c.Next()
	status := c.Writer.Status()
	statusStr := strconv.Itoa(status)
	httpRequestsTotal.WithLabelValues(method, path, statusStr, instance).Inc()
	timer.ObserveDuration()
}
