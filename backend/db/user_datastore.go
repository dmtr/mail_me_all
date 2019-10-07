package db

import (
	"fmt"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	UniqueViolationErr pq.ErrorCode = "23505"
)

type DbError struct {
	PqError *pq.Error
	Err     error
}

func (e *DbError) Error() string {
	if e.PqError != nil {
		return fmt.Sprintf("Got Database error with code: %s message: %s detail: %s and constraint: %s",
			e.PqError.Code, e.PqError.Message, e.PqError.Detail, e.PqError.Constraint)
	}
	return e.Err.Error()
}

func getPqError(err error) *pq.Error {
	if err, ok := err.(*pq.Error); ok {
		return err
	} else {
		return nil
	}
}

func getDbError(err error) error {
	if err != nil {
		return &DbError{
			PqError: getPqError(err),
			Err:     err,
		}
	}
	return err
}

func rollbackOnError(tx *sqlx.Tx, err error) {
	if err != nil {
		log.Info("Rollback")
		tx.Rollback()
	}
}

type UserDatastore struct {
	DB *sqlx.DB
}

func NewUserDatastore(db *sqlx.DB) *UserDatastore {
	return &UserDatastore{DB: db}
}

func (d *UserDatastore) CreateUser(user *models.User) error {
	log.Debugf("Going to insert user %v", user)

	tx := d.DB.MustBegin()
	var err error

	defer func() {
		rollbackOnError(tx, err)
	}()

	_, err = tx.NamedExec("INSERT INTO user_account (name, fb_id, fb_token) VALUES (:name, :fb_id, :fb_token)", user)
	err = getDbError(err)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" inserting user: %v", user))
		return err
	}

	err = getDbError(tx.Commit())
	return err
}
