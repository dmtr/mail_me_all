package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })
	router.GET("/oauth/fb", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!!")
		log.Infof("context %+v", c)
	})
}
