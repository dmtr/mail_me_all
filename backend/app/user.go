package app

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
)

// CreateUser - add new user
func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	t, exists := c.Get("Tx")
	if !exists {
		log.Error("No transaction in context!")
		return
	}
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		log.Error("No transaction in context!")
	}

	_, err := tx.NamedExec("INSERT INTO user_account (name, fb_id, fb_token) VALUES (:name, :fb_id, :fb_token)", &user)
	if err != nil {
		log.Errorf("Got error %v", err)
		panic(err)
	}
	log.Debugf("New user %v", user)
	c.JSON(http.StatusCreated, gin.H{"name": user.Name})
}
