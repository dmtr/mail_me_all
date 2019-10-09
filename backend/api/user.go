package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type fbuser struct {
	ID    string `json:"fbid" binding:"required"`
	Token string `json:"fbtoken" binding:"required"`
}

func setSessionCookie(c *gin.Context, conf *config.Config, userID string) {
	s := sessions.Default(c)
	uid := s.Get("userid")
	log.Debugf("Got from session %s", uid)
	if uid == nil {
		s.Options(sessions.Options{
			Path:     conf.Path,
			Domain:   conf.Domain,
			MaxAge:   conf.MaxAge,
			Secure:   conf.Secure,
			HttpOnly: conf.HttpOnly,
		})
		s.Set("userid", userID)
		err := s.Save()
		if err != nil {
			log.Errorf("Can not start session, error %s", err)
		}
	}
}

// SignInFB - sign in with Facebook
func SignInFB(conf *config.Config, usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user fbuser
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Debugf("Id %s", user.ID)
		err := usecases.User.SignInFB(user.ID, user.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		} else {
			log.Debugf("Signed in with Facebook %s", user.ID)
			setSessionCookie(c, conf, user.ID)
			c.JSON(http.StatusOK, gin.H{"id": user.ID})
		}
	}
}
