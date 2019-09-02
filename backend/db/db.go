package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	maxIdleConns = 10
	maxOpenConns = 16
)

// ConnectDb - try to connect to the database
func ConnectDb(DSN string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := sqlx.Connect("postgres", DSN)
			if err == nil {
				db.SetMaxIdleConns(maxIdleConns)
				db.SetMaxOpenConns(maxOpenConns)
				return db, nil
			}
			log.Errorf("failed to connect to db %s %s", DSN, err)
		}
	}
}
