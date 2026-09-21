package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const placeholderPage = `enchantech backend is running.

The site's front end is not part of this repository. Point APP_CONTENT_DIR at a
directory containing templates/home.html and assets/, or place one at ./content,
and this route will render the page instead.
`

// registerRoutes keeps all HTTP route registration in one place.
func registerRoutes(router *gin.Engine, hasPage bool) {
	router.GET("/", func(context *gin.Context) {
		if !hasPage {
			context.String(http.StatusOK, placeholderPage)

			return
		}

		context.HTML(http.StatusOK, "home.html", newPageData(siteProfile))
	})
}
