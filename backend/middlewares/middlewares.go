package middlewares

import (
	"net/http"

	"github.com/dmtr/mail_me_all/backend/errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/sessions"
)

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

// SessionlMiddleware adds user id to the context if session exists, otherwise returns 401
func SessionlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUserID(c)
		if uid != "" {
			c.Set("UserID", uid)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": errors.AuthRequired})
		}
	}
}

// TestSessionMiddleware must be used in tests
func TestSessionMiddleware(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("UserID", userID)
	}
}

// TransactionlMiddleware - use for requests that need db transaction
func TransactionlMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.MustBegin()
		c.Set("Tx", tx)
		defer func() {
			if err := recover(); err != nil {
				log.Errorln("Panic, transaction rollback")
				tx.Rollback()
				panic(err)
			} else {
				log.Debugln("Commiting transaction")
				tx.Commit()
			}
		}()
		c.Next()
	}
}

// TestTransactionlMiddleware - must be used in tests
func TestTransactionlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := &sqlx.Tx{}
		c.Set("Tx", tx)
		c.Next()
	}
}
