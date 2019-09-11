package api

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// CreateUser - add new user
func CreateUser(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)
		err := usecases.User.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		} else {
			log.Debugf("New user %v", user)
			c.JSON(http.StatusCreated, gin.H{"name": user.Name})
		}
	}
}
