package server

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	return gin.Default()
}
