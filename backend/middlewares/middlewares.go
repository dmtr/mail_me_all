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
