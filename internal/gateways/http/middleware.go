package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}
		c.Next()
	}
}

func RequireJSONAccept() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept")
		if !strings.HasPrefix(accept, "application/json") {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
		c.Next()
	}
}
