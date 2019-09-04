package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// TransactionlMiddleware - use for requests that need db transaction
func TransactionlMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.MustBegin()
		c.Set("Tx", tx)
		c.Next()
		if c.IsAborted() {
			log.Errorln("Context is aborted, transaction rollback")
			tx.Rollback()
		} else {
			log.Debugln("Commiting transaction")
			tx.Commit()
		}
	}
}
