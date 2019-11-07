package db

import (
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
