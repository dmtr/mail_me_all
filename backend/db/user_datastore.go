package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	uniqueViolationErr pq.ErrorCode = "23505"
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

func (e *DbError) HasNoRows() bool {
	return e.Err == sql.ErrNoRows
}

func (e *DbError) IsUniqueViolationError() bool {
	return e.PqError.Code == uniqueViolationErr
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

type subscriptionUser struct {
	ID            uint   `db:"id"`
	Name          string `db:"name"`
	TwitterID     string `db:"twitter_id"`
	ProfileIMGURL string `db:"profile_image_url"`
	ScreenName    string `db:"screen_name"`
}

func (s *subscriptionUser) insertSubscriptionUser(tx *sqlx.Tx) error {
	res, err := tx.NamedQuery("INSERT INTO subscription_user (twitter_id, name, profile_image_url, screen_name) VALUES (:twitter_id, :name, :profile_image_url, :screen_name) RETURNING id", *s)
	if err != nil {
		return err
	}

	var id uint
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			log.Errorf("Scan error: %s", err)
			return err
		}
	}
	s.ID = id
	return err
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
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		res, err := tx.NamedQuery("INSERT INTO user_account (name, email) VALUES (:name, :email) RETURNING id", user)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" inserting user: %s", user))
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

func (d *UserDatastore) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	log.Debugf("Going to update user %s", user.ID)
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		_, err := tx.NamedExec("UPDATE user_account SET name=:name, email=:email WHERE id = :id", user)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" uptating user: %s", user))
			return models.User{}, err
		}
		return user, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.User{}, err
	}

	r, _ := res.(models.User)
	return r, err
}

func (d *UserDatastore) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	tx := getTransaction(ctx)

	f := func(tx *sqlx.Tx) (models.Model, error) {
		var user models.User
		err := tx.Get(&user, "SELECT id, name, email FROM user_account WHERE id=$1", userID)
		return user, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.User{}, err
	}

	r, _ := res.(models.User)
	return r, err
}

func (d *UserDatastore) GetTwitterUserByID(ctx context.Context, twitterUserID string) (models.TwitterUser, error) {
	tx := getTransaction(ctx)

	f := func(tx *sqlx.Tx) (models.Model, error) {
		var user models.TwitterUser
		err := tx.Get(
			&user, "SELECT user_id, social_account_id, access_token, token_secret, profile_image_url FROM tw_account WHERE social_account_id=$1", twitterUserID)
		return user, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.TwitterUser{}, err
	}

	r, _ := res.(models.TwitterUser)
	return r, err
}

func (d *UserDatastore) InsertTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		_, err := tx.NamedExec("INSERT INTO tw_account (user_id, social_account_id, access_token, token_secret, profile_image_url) VALUES (:user_id, :social_account_id, :access_token, :token_secret, :profile_image_url)", twitterUser)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" inserting twitterUser: %s", twitterUser))
			return models.TwitterUser{}, err
		}
		return twitterUser, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.TwitterUser{}, err
	}

	r, _ := res.(models.TwitterUser)
	return r, err
}

func (d *UserDatastore) UpdateTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	log.Debugf("Going to update twitterUser for user %s", twitterUser.UserID)
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		_, err := tx.NamedExec("UPDATE tw_account SET access_token=:access_token, token_secret=:token_secret, profile_image_url=:profile_image_url WHERE user_id = :user_id", twitterUser)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" uptating twitterUser: %s", twitterUser))
			return models.TwitterUser{}, err
		}
		return twitterUser, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.TwitterUser{}, err
	}

	r, _ := res.(models.TwitterUser)
	return r, err
}

func (d *UserDatastore) GetTwitterUser(ctx context.Context, userID uuid.UUID) (models.TwitterUser, error) {
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		var twitterUser models.TwitterUser
		err := tx.Get(&twitterUser, "SELECT user_id, social_account_id, access_token, token_secret, profile_image_url FROM tw_account WHERE user_id = $1", userID)
		return twitterUser, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.TwitterUser{}, err
	}

	r, _ := res.(models.TwitterUser)
	return r, err
}

func (d *UserDatastore) InsertSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	tx := getTransaction(ctx)
	f := func(tx *sqlx.Tx) (models.Model, error) {
		res, err := tx.NamedQuery("INSERT INTO subscription (user_id, title, email, day) VALUES (:user_id, :title, :email, :day) RETURNING id", subscription)
		if err != nil {
			log.Error(err.Error() + fmt.Sprintf(" inserting subscription: %s", subscription))
			return models.Subscription{}, err
		}

		var id string
		for res.Next() {
			err = res.Scan(&id)
			if err != nil {
				log.Errorf("Scan error: %s", err)
				return subscription, err
			}
		}
		subscription.ID, err = uuid.Parse(id)
		if err != nil {
			log.Errorf("Can not parse subscription id %s", id)
			return subscription, err
		}

		for _, u := range subscription.UserList {
			su := &subscriptionUser{
				TwitterID:     u.TwitterID,
				Name:          u.Name,
				ProfileIMGURL: u.ProfileIMGURL,
				ScreenName:    u.ScreenName,
			}

			_, err = tx.Exec("SAVEPOINT save1")
			if err != nil {
				break
			}

			err = su.insertSubscriptionUser(tx)

			if err != nil {
				e := getDbError(err).(*DbError)
				if e.IsUniqueViolationError() {
					_, err = tx.Exec("ROLLBACK TO SAVEPOINT save1")
					if err != nil {
						break
					}

					err = nil
					continue
				} else {
					break
				}
			}
			m2m := struct {
				Subscription_id uuid.UUID `db:"subscription_id"`
				User_id         uint      `db:"user_id"`
			}{
				subscription.ID,
				su.ID,
			}
			_, err = tx.NamedExec("INSERT INTO subscription_user_m2m (subscription_id, user_id) VALUES(:subscription_id, :user_id)", m2m)
			if err != nil {
				break
			}
		}

		return subscription, err
	}

	res, err := d.execQuery(tx, f)
	if err != nil {
		return models.Subscription{}, err
	}

	r, _ := res.(models.Subscription)
	return r, err
}
