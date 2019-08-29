package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// User - represents user
type User struct {
	Name    string `json:"name" binding:"required"`
	FbID    string `json:"fbid"`
	FbToken string `json:"fbtoken"`
}

// CreateUser - add new user
func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	log.Infof("New user %v", user)
	c.JSON(http.StatusOK, gin.H{"name": user.Name})
}
