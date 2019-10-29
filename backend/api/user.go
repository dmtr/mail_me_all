package api

import (
	"context"
	"fmt"
	"net/http"

	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
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

type appUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SignedIn bool   `json:"signedIn"`
}

func adaptUser(user models.User, signedIn bool) appUser {
	return appUser{
		ID:       user.ID.String(),
		Name:     user.Name,
		SignedIn: signedIn,
	}

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

func setSessionCookie(c *gin.Context, conf *config.Config, userID string) error {
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
	return err
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

func ProcessTwitterCallback(conf *config.Config, oauth1Config *oauth1.Config, usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {

		success := func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			twitterUser, err := twitter.UserFromContext(ctx)
			log.Debugf("Twitter user: %v", twitterUser)

			if err != nil {
				log.Errorf("Callback error %s", err.Error())
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			accessToken, accessSecret, err := oauth1Login.AccessTokenFromContext(ctx)
			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			tx, err := getTransaction(c)
			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			contxt := context.WithValue(context.Background(), "Tx", tx)

			user, err := usecases.User.SignInWithTwitter(
				contxt, twitterUser.IDStr, twitterUser.Name, twitterUser.Email, accessToken, accessSecret)

			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			setSessionCookie(c, conf, user.ID.String())
			c.Redirect(http.StatusFound, "/")
		}

		h := twitter.CallbackHandler(oauth1Config, http.HandlerFunc(success), nil)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// GetUser - get user id from session cookie and check if user is valid
func GetUser(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u appUser
		uid := getUserID(c)
		if uid == "" {
			c.JSON(http.StatusNotFound, u)
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
				c.JSON(http.StatusOK, adaptUser(user, true))
			}
		}
	}
}
