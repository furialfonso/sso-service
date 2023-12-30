package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IPingHandler interface {
	Ping(c *gin.Context)
}

type pingHandler struct{}

func NewHandlerPing() IPingHandler {
	return &pingHandler{}
}

func (h *pingHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
