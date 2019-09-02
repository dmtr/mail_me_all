package app

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
)

// GetRouter - returns router
func GetRouter(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, db)
	return router
}

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine, db *sqlx.DB) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })

	router.GET("/oauth/fb", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!!")
		log.Debugf("context %+v", c)
	})

	router.POST("/api/users", middlewares.TransactionlMiddleware(db), CreateUser)
}
