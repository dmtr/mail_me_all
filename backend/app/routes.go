package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// GetRouter - returns router
func GetRouter() *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router)
	return router
}

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })

	router.GET("/oauth/fb", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!!")
		log.Debugf("context %+v", c)
	})

	router.POST("/api/users", CreateUser)
}
