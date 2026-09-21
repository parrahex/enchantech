package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerRoutes keeps all HTTP route registration in one place.
func registerRoutes(router *gin.Engine) {
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "enchantech backend is running.\n")
	})
}
