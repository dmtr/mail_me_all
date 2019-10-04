package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type fbuser struct {
	Id    string `json:"fbid" binding:"required"`
	Token string `json:"fbtoken" binding:"required"`
}

// SignInFB - sign in with Facebook
func SignInFB(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user fbuser
		c.BindJSON(&user)
		err := usecases.User.SignInFB(user.Id, user.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		} else {
			log.Debugf("Signed in with Facebook %s", user.Id)
			c.JSON(http.StatusCreated, gin.H{"id": user.Id})
		}
	}
}
