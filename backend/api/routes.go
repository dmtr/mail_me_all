package api

import (
	"net/http"

	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/middlewares"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
func GetRouter(conf *config.Config, db *sqlx.DB, usecases *models.UseCases) *gin.Engine {
	router := gin.Default()
	sessionStore := getSessionStore(conf.AuthKey, conf.EncryptKey)
	router.Use(sessions.Sessions("session", sessionStore))
	RegisterRoutes(router, conf, db, usecases, conf.Testing)
	return router
}

//RegisterRoutes setups routes
func RegisterRoutes(router *gin.Engine, conf *config.Config, db *sqlx.DB, usecases *models.UseCases, testing bool) {
	router.GET("/healthcheck", func(c *gin.Context) { c.String(http.StatusOK, "Ok") })

	oauth1Config := &oauth1.Config{
		ConsumerKey:    conf.TwConsumerKey,
		ConsumerSecret: conf.TwConsumerSecret,
		CallbackURL:    conf.TwCallbackURL,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}

	if testing { // unit tests
		api := router.Group("/api", middlewares.TestSessionMiddleware(testUserID))
		api.GET("/user", middlewares.TestTransactionlMiddleware(), GetUser(usecases))
		api.GET("/twitter-users", SearchTwitterUsers(usecases))
	} else {
		router.GET("/oauth/tw/signin", gin.WrapH(twitter.LoginHandler(oauth1Config, nil)))
		router.GET("/oauth/tw/callback", middlewares.TransactionlMiddleware(db), ProcessTwitterCallback(conf, oauth1Config, usecases))

		api := router.Group("/api", middlewares.SessionMiddleware())
		api.GET("/user", middlewares.TransactionlMiddleware(db), GetUser(usecases))
		api.GET("/twitter-users", SearchTwitterUsers(usecases))
	}
}
