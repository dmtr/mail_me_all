package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })
}
