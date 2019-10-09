package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func getSessionStore(authKey string, encryptKey string) cookie.Store {
	var keys [][]byte
	keys = append(keys, []byte(authKey))

	if encryptKey != "" {
		keys = append(keys, []byte(encryptKey))
	}

	return cookie.NewStore(keys...)
}

// GetRouter - returns router
func GetRouter(conf *config.Config, usecases *models.UseCases) *gin.Engine {
	router := gin.Default()
	sessionStore := getSessionStore(conf.AuthKey, conf.EncryptKey)
	router.Use(sessions.Sessions("session", sessionStore))
	RegisterRoutes(router, conf, usecases)
	return router
}

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine, conf *config.Config, usecases *models.UseCases) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })

	router.GET("/oauth/fb", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!!")
		log.Debugf("context %+v", c)
	})

	router.POST("/api/signin/fb", SignInFB(conf, usecases))
}
