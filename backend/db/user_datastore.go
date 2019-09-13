package db

import (
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	uniqueViolationErr pq.ErrorCode = "23505"
)

type UserDatastore struct {
	DB *sqlx.DB
}

func NewUserDatastore(db *sqlx.DB) *UserDatastore {
	return &UserDatastore{DB: db}
}

func (d *UserDatastore) CreateUser(user *models.User) error {
	log.Debugf("Going to insert user %v", user)
	tx := d.DB.MustBegin()
	_, err := tx.NamedExec("INSERT INTO user_account (name, fb_id, fb_token) VALUES (:name, :fb_id, :fb_token)", user)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Errorf("Got error with code: %s message: %s detail: %s and constraint: %s inserting user: %v", err.Code, err.Message, err.Detail, err.Constraint, user)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Got error commiting transaction %v", err)
	}

	return err
}
