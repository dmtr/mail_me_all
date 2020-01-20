package db

import (
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type subscription struct {
	SubscriptionID uuid.UUID `db:"subscription_id"`
	Title          string    `db:"title"`
	Email          string    `db:"email"`
	Day            string    `db:"day"`
	UserID         uuid.UUID `db:"user_id"`
	IgnoreRT       bool      `db:"ignore_rt"`
	IgnoreReplies  bool      `db:"ignore_replies"`
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

func getSubscriptionUser(tx *sqlx.Tx, twitterID string) (subscriptionUser, error) {
	var user subscriptionUser
	err := tx.Get(&user, "SELECT id, name, twitter_id, profile_image_url, screen_name FROM subscription_user WHERE twitter_id=$1", twitterID)
	return user, err
}

func insertUserList(tx *sqlx.Tx, userList models.UserList, subscriptionID uuid.UUID) error {
	var err error
	for _, u := range userList {
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

				*su, err = getSubscriptionUser(tx, su.TwitterID)
				if err != nil {
					break
				}

				err = nil

			} else {
				break
			}
		}
		m2m := struct {
			Subscription_id uuid.UUID `db:"subscription_id"`
			User_id         uint      `db:"user_id"`
		}{
			subscriptionID,
			su.ID,
		}
		_, err = tx.NamedExec("INSERT INTO subscription_user_m2m (subscription_id, user_id) VALUES(:subscription_id, :user_id)", m2m)
		if err != nil {
			break
		}
	}
	return err
}
