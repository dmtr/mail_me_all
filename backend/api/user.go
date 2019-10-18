package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/errors"
	"github.com/dmtr/mail_me_all/backend/models"
	useCases "github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
)

type fbuser struct {
	ID    string `json:"fbid" binding:"required"`
	Token string `json:"fbtoken" binding:"required"`
}

type appUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SignedIn bool   `json:"signedIn"`
}

func getUserID(c *gin.Context) string {
	s := sessions.Default(c)
	uid := s.Get("userid")
	if uid == nil {
		return ""
	}
	u, ok := uid.(string)
	if !ok {
		log.Warningf("Can not convert userid to string %v", uid)
		return ""
	}
	return u
}

func setSessionCookie(c *gin.Context, conf *config.Config, userID string) {
	s := sessions.Default(c)
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

func getTransaction(c *gin.Context) (*sqlx.Tx, error) {
	t, exists := c.Get("Tx")
	if !exists {
		return nil, fmt.Errorf("No transaction in context!")
	}
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return nil, fmt.Errorf("Wrong transaction type!")
	}
	return tx, nil
}

// SignInFB - sign in with Facebook
func SignInFB(conf *config.Config, usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user fbuser
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest})
			return
		}
		log.Debugf("Id %s", user.ID)

		tx, err := getTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		ctx := context.WithValue(context.Background(), "Tx", tx)
		u, err := usecases.User.SignInFB(ctx, user.ID, user.Token)
		if err != nil {
			log.Errorf("Can not sign in %s", err)
			e, _ := err.(*useCases.UseCaseError)
			c.JSON(http.StatusInternalServerError, gin.H{"code": e.Code()})
		} else {
			log.Debugf("Signed in with Facebook %s", user.ID)
			setSessionCookie(c, conf, u.ID.String())
			c.JSON(http.StatusOK, gin.H{"id": u.ID})
		}
	}
}

// GetUser - get user id from session cookie and check if user is valid
func GetUser(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u appUser
		uid := getUserID(c)
		if uid == "" {
			c.JSON(http.StatusOK, u)
		} else {
			tx, err := getTransaction(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
				return
			}

			ctx := context.WithValue(context.Background(), "Tx", tx)

			userID, err := uuid.Parse(uid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
				return
			}

			user, err := usecases.User.GetUserByID(ctx, userID)
			if err != nil {
				log.Errorf("Can not get user, got error %s", err)
				e, _ := err.(*useCases.UseCaseError)
				c.JSON(http.StatusInternalServerError, gin.H{"code": e.Code()})
			} else {
				u.ID = user.ID.String()
				u.Name = user.Name
				u.SignedIn = true
				c.JSON(http.StatusOK, u)
			}
		}
	}
}
