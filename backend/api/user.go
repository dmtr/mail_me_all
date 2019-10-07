package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type fbuser struct {
	ID    string `json:"fbid" binding:"required"`
	Token string `json:"fbtoken" binding:"required"`
}

// SignInFB - sign in with Facebook
func SignInFB(usecases *models.UseCases) gin.HandlerFunc {
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
			c.JSON(http.StatusOK, gin.H{"id": user.ID})
		}
	}
}
