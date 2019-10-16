package db

import (
	"context"
	"fmt"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	UniqueViolationErr pq.ErrorCode = "23505"
)

type queryFunc func(tx *sqlx.Tx) (models.Model, error)

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

func getTransaction(ctx context.Context) *sqlx.Tx {
	t := ctx.Value("Tx")
	if t == nil {
		return nil
	}

	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}

func (d *UserDatastore) execQuery(tx *sqlx.Tx, f queryFunc) (models.Model, error) {
	var err error
	beginTx := false

	if tx == nil {
		beginTx = true
		tx = d.DB.MustBegin()

		defer func() {
			rollbackOnError(tx, err)
		}()
	}

	var res models.Model
	res, err = f(tx)

	err = getDbError(err)

	if beginTx {
		err = getDbError(tx.Commit())
	}
	return res, err
}

func (d *UserDatastore) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	log.Debugf("Going to insert user %s", user.FbID)
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		res, err := tx.NamedQuery("INSERT INTO user_account (name, fb_id, email) VALUES (:name, :fb_id, :email) RETURNING id", user)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" inserting user: %v", user))
			return user, err
		}

		var userId string
		for res.Next() {
			err = res.Scan(&userId)
			if err != nil {
				log.Errorf("Scan error: %s", err)
				return user, err
			}
		}
		log.Debugf("Got user id %s", userId)
		user.ID, err = uuid.Parse(userId)
		if err != nil {
			log.Errorf("Can not parse user id %s", userId)
			return user, err
		}
		return user, nil
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.User{}, err
	}

	r, _ := res.(models.User)
	return r, err
}

func (d *UserDatastore) InsertToken(ctx context.Context, token models.Token) (models.Token, error) {
	log.Debugf("Going to insert token for user %s", token.UserID)
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		_, err := tx.NamedExec("INSERT INTO token (user_id, fb_token, expires_at) VALUES (:user_id, :fb_token, :expires_at)", token)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" inserting token: %v", token))
			return token, err
		}
		return models.Token{}, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.Token{}, err
	}

	r, _ := res.(models.Token)
	return r, err
}

func (d *UserDatastore) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	tx := getTransaction(ctx)

	f := func(tx *sqlx.Tx) (models.Model, error) {
		var user models.User
		err := tx.Get(&user, "SELECT id, name, email, fb_id FROM user_account WHERE id=$1", userID)
		return user, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.User{}, err
	}

	r, _ := res.(models.User)
	return r, err
}
