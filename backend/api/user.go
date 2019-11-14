package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

type twitterUser struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ScreenName    string `json:"screen_name"`
	ProfileIMGURL string `json:"profile_image_url"`
}

type subscription struct {
	ID       string        `json:"id"`
	Title    string        `json:"title" binding:"required"`
	Email    string        `json:"email" binding:"required"`
	Day      string        `json:"day" binding:"required"`
	UserList []twitterUser `json:"userList" binding:"required"`
}

func adaptUser(user models.User, signedIn bool) appUser {
	return appUser{
		ID:       user.ID.String(),
		Name:     user.Name,
		SignedIn: signedIn,
	}

}

func adaptTwitterUserSearchResult(user models.TwitterUserSearchResult) twitterUser {
	return twitterUser{
		ID:            user.TwitterID,
		Name:          user.Name,
		ScreenName:    user.ScreenName,
		ProfileIMGURL: user.ProfileIMGURL,
	}
}

func adaptSubscription(s models.Subscription) subscription {
	subcr := subscription{
		ID:    s.ID.String(),
		Title: s.Title,
		Email: s.Email,
		Day:   s.Day,
	}

	for _, u := range s.UserList {
		subcr.UserList = append(subcr.UserList, adaptTwitterUserSearchResult(u))
	}
	return subcr
}

func getUserID(c *gin.Context) string {
	uid, exists := c.Get("UserID")
	if !exists {
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

func getContextWithTransaction(c *gin.Context) (context.Context, error) {
	tx, err := getTransaction(c)
	if err != nil {
		return nil, err
	}

	return context.WithValue(context.Background(), "Tx", tx), nil
}

func processTwitterCallback(conf *config.Config, oauth1Config *oauth1.Config, usecases *models.UseCases) gin.HandlerFunc {
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
				log.Errorf("Error, no tokens: %s", err.Error())
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			contxt, err := getContextWithTransaction(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
				return
			}

			user, err := usecases.User.SignInWithTwitter(
				contxt, twitterUser.IDStr, twitterUser.Name, twitterUser.Email, twitterUser.ScreenName, accessToken, accessSecret)

			if err != nil {
				log.Errorf("Can not sign in with Twitter: %s", err.Error())
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
func getUser(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u appUser
		uid := getUserID(c)
		if uid == "" {
			c.JSON(http.StatusNotFound, u)
			return
		}

		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		ctx, err := getContextWithTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		user, err := usecases.User.GetUserByID(ctx, userID)
		if err != nil {
			log.Errorf("Can not get user, got error %s", err)
			e, _ := err.(*useCases.UseCaseError)
			status := http.StatusInternalServerError
			if e.Code() == errors.NotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"code": e.Code()})
			return
		}

		c.JSON(http.StatusOK, adaptUser(user, true))
	}
}

func searchTwitterUsers(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		if uid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest})
			return
		}

		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest})
			return
		}

		users, err := usecases.User.SearchTwitterUsers(context.Background(), userID, query)
		if err != nil {
			log.Errorf("Got error searching users %s", err)
			e, _ := err.(*useCases.UseCaseError)
			c.JSON(http.StatusInternalServerError, gin.H{"code": e.Code()})
			return
		}

		res := make([]twitterUser, 0, len(users))
		for _, user := range users {
			res = append(res, adaptTwitterUserSearchResult(user))
		}

		c.JSON(http.StatusOK, gin.H{"users": res})
	}
}

func addSubscription(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		newSubscription, err := getSubscription(c, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		ctx, err := getContextWithTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		newSubscription, err = usecases.User.AddSubscription(ctx, newSubscription)

		if err != nil {
			log.Errorf("Can not add subscription %s, got error %s", newSubscription, err)
			e, _ := err.(*useCases.UseCaseError)
			status := http.StatusInternalServerError
			if e.Code() == errors.NotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"code": e.Code()})
			return
		}

		c.JSON(http.StatusOK, adaptSubscription(newSubscription))
	}
}

func getSubscriptions(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		ctx, err := getContextWithTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		subscriptions, err := usecases.User.GetSubscriptions(ctx, userID)
		if err != nil {
			log.Errorf("Can not find subscriptions for user %s, got error %s", userID, err)
			e, _ := err.(*useCases.UseCaseError)
			status := http.StatusInternalServerError
			if e.Code() == errors.NotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"code": e.Code()})
			return
		}

		res := make([]subscription, 0, len(subscriptions))
		for _, s := range subscriptions {
			res = append(res, adaptSubscription(s))
		}

		c.JSON(http.StatusOK, gin.H{"subscriptions": res})
	}
}

func getSubscription(c *gin.Context, userID uuid.UUID) (models.Subscription, error) {
	var s subscription
	if err := c.ShouldBindJSON(&s); err != nil {
		return models.Subscription{}, err
	}

	id, _ := uuid.Parse(s.ID)
	newSubscription := models.Subscription{
		ID:     id,
		UserID: userID,
		Title:  s.Title,
		Email:  s.Email,
		Day:    strings.ToLower(s.Day),
	}

	for _, u := range s.UserList {
		newSubscription.UserList = append(
			newSubscription.UserList, models.TwitterUserSearchResult{
				TwitterID:     u.ID,
				Name:          u.Name,
				ProfileIMGURL: u.ProfileIMGURL,
				ScreenName:    u.ScreenName})
	}
	return newSubscription, nil
}

func updateSubscription(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		ctx, err := getContextWithTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		updatedSubscription, err := getSubscription(c, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		updatedSubscription, err = usecases.User.UpdateSubscription(ctx, userID, updatedSubscription)
		if err != nil {
			log.Errorf("Can not add subscription %s, got error %s", updatedSubscription, err)
			e, _ := err.(*useCases.UseCaseError)
			status := http.StatusInternalServerError
			if e.Code() == errors.NotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"code": e.Code()})
			return
		}

		c.JSON(http.StatusOK, adaptSubscription(updatedSubscription))
	}
}

func deleteSubscription(usecases *models.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		userID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": errors.BadRequest, "message": err.Error()})
			return
		}

		subscriptionID, err := uuid.Parse(c.Param("id"))

		ctx, err := getContextWithTransaction(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": errors.ServerError})
			return
		}

		err = usecases.User.DeleteSubscription(ctx, userID, subscriptionID)
		if err != nil {
			log.Errorf("Can not delete subscription %s, got error %s", subscriptionID, err)
			e, _ := err.(*useCases.UseCaseError)
			status := http.StatusInternalServerError

			if e.Code() == errors.NotFound {
				status = http.StatusNotFound
			} else if e.Code() == errors.AuthRequired {
				status = http.StatusUnauthorized
			}

			c.JSON(status, gin.H{"code": e.Code()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
