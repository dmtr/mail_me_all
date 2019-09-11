package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// GetRouter - returns router
func GetRouter(usecases *models.UseCases) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, usecases)
	return router
}

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine, usecases *models.UseCases) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })

	router.GET("/oauth/fb", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!!")
		log.Debugf("context %+v", c)
	})

	router.POST("/api/users", CreateUser(usecases))
}
