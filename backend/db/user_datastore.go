package db

import (
	"context"
	"fmt"

	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type currentTx struct {
	tx      *sqlx.Tx
	beginTx bool
	err     *error
}

func getTransactionFromContext(ctx context.Context) *sqlx.Tx {
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

func getTransaction(ctx context.Context, db *sqlx.DB, err *error) currentTx {
	var beginTx bool
	tx := getTransactionFromContext(ctx)
	if tx == nil {
		tx = db.MustBegin()
		beginTx = true
	}
	return currentTx{tx, beginTx, err}
}

func (t *currentTx) commitOrRollback() {
	if t.beginTx {
		if *t.err != nil {
			log.Info("Rollback")
			t.tx.Rollback()
		} else {
			e := t.tx.Commit()
			if e != nil {
				log.Errorf("Error committing transaction %s", e)
			}
		}
	}
}

func (t *currentTx) getError() error {
	return getDbError(*t.err)
}

type UserDatastore struct {
	DB *sqlx.DB
}

func NewUserDatastore(db *sqlx.DB) *UserDatastore {
	return &UserDatastore{DB: db}
}

func (d *UserDatastore) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	res, err := t.tx.NamedQuery("INSERT INTO user_account (name, email) VALUES (:name, :email) RETURNING id", user)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" inserting user: %s", user))
		return user, t.getError()
	}

	var userId string
	for res.Next() {
		err = res.Scan(&userId)
		if err != nil {
			log.Errorf("Scan error: %s", err)
			return user, t.getError()
		}
	}
	log.Debugf("Got user id %s", userId)
	user.ID, err = uuid.Parse(userId)
	if err != nil {
		log.Errorf("Can not parse user id %s", userId)
		return user, t.getError()
	}
	return user, nil
}

func (d *UserDatastore) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	log.Debugf("Going to update user %s", user.ID)
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.NamedExec("UPDATE user_account SET name=:name, email=:email WHERE id = :id", user)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" uptating user: %s", user))
		return models.User{}, t.getError()
	}
	return user, t.getError()

}

func (d *UserDatastore) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var user models.User
	err = t.tx.Get(&user, "SELECT id, name, email FROM user_account WHERE id=$1", userID)
	return user, t.getError()

}

func (d *UserDatastore) GetTwitterUserByID(ctx context.Context, twitterUserID string) (models.TwitterUser, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var user models.TwitterUser
	err = t.tx.Get(
		&user, "SELECT user_id, social_account_id, access_token, token_secret, profile_image_url FROM tw_account WHERE social_account_id=$1", twitterUserID)
	return user, t.getError()

}

func (d *UserDatastore) InsertTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.NamedExec("INSERT INTO tw_account (user_id, social_account_id, access_token, token_secret, profile_image_url) VALUES (:user_id, :social_account_id, :access_token, :token_secret, :profile_image_url)", twitterUser)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" inserting twitterUser: %s", twitterUser))
		return models.TwitterUser{}, t.getError()
	}
	return twitterUser, t.getError()

}

func (d *UserDatastore) UpdateTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	log.Debugf("Going to update twitterUser for user %s", twitterUser.UserID)
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()
	_, err = t.tx.NamedExec("UPDATE tw_account SET access_token=:access_token, token_secret=:token_secret, profile_image_url=:profile_image_url WHERE user_id = :user_id", twitterUser)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" uptating twitterUser: %s", twitterUser))
		return models.TwitterUser{}, t.getError()
	}
	return twitterUser, t.getError()
}

func (d *UserDatastore) GetTwitterUser(ctx context.Context, userID uuid.UUID) (models.TwitterUser, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var twitterUser models.TwitterUser
	err = t.tx.Get(&twitterUser, "SELECT user_id, social_account_id, access_token, token_secret, profile_image_url FROM tw_account WHERE user_id = $1", userID)
	return twitterUser, t.getError()

}

func (d *UserDatastore) InsertSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	tx := t.tx
	res, err := tx.NamedQuery("INSERT INTO subscription (user_id, title, email, day) VALUES (:user_id, :title, :email, :day) RETURNING id", subscription)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" inserting subscription: %s", subscription))
		return models.Subscription{}, t.getError()
	}

	var id string
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			log.Errorf("Scan error: %s", err)
			return subscription, t.getError()
		}
	}
	subscription.ID, err = uuid.Parse(id)
	if err != nil {
		log.Errorf("Can not parse subscription id %s", id)
		return subscription, t.getError()
	}

	err = insertUserList(tx, subscription.UserList, subscription.ID)

	return subscription, t.getError()
}

func (d *UserDatastore) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	type subscriptionRow struct {
		subscription
		subscriptionUser
	}

	rows, err := t.tx.Queryx(
		"SELECT s.id AS subscription_id, s.user_id, s.title, s.email, s.day, u.id, u.name, u.twitter_id, u.screen_name, u.profile_image_url "+
			"FROM subscription s "+
			"INNER JOIN subscription_user_m2m m2m ON m2m.subscription_id = s.id "+
			"INNER JOIN subscription_user u ON u.id = m2m.user_id "+
			"WHERE s.user_id = $1 "+
			"ORDER BY s.updated_at DESC", userID)

	if err != nil {
		return []models.Subscription{}, t.getError()
	}

	processed := make(map[uuid.UUID]models.Subscription)
	processedKeys := make([]uuid.UUID, 0)

	for rows.Next() {
		var row subscriptionRow
		err = rows.StructScan(&row)

		s := models.Subscription{
			ID:     row.SubscriptionID,
			Title:  row.Title,
			Email:  row.Email,
			Day:    row.Day,
			UserID: row.UserID,
		}
		u := models.TwitterUserSearchResult{
			TwitterID:     row.TwitterID,
			Name:          row.Name,
			ProfileIMGURL: row.ProfileIMGURL,
			ScreenName:    row.ScreenName,
		}
		processedSubscription, ok := processed[s.ID]
		if ok {
			processedSubscription.UserList = append(processedSubscription.UserList, u)
			processed[s.ID] = processedSubscription
		} else {
			s.UserList = append(s.UserList, u)
			processed[s.ID] = s
			processedKeys = append(processedKeys, s.ID)
		}
	}

	res := make([]models.Subscription, 0, len(processed))
	for _, k := range processedKeys {
		res = append(res, processed[k])
	}

	return res, t.getError()
}

func (d *UserDatastore) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (models.Subscription, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var subscription models.Subscription

	err = t.tx.Get(&subscription, "SELECT id, user_id, title, email, day FROM subscription WHERE id=$1", subscriptionID)

	if err != nil {
		return subscription, t.getError()
	}

	rows, err := t.tx.Queryx("SELECT u.id, u.name, u.twitter_id, u.profile_image_url, u.screen_name FROM subscription_user u "+
		"INNER JOIN subscription_user_m2m m ON u.id = m.user_id "+
		"WHERE m.subscription_id=$1", subscriptionID)

	if err != nil {
		return subscription, t.getError()
	}

	for rows.Next() {
		var row subscriptionUser
		err = rows.StructScan(&row)

		u := models.TwitterUserSearchResult{
			TwitterID:     row.TwitterID,
			Name:          row.Name,
			ProfileIMGURL: row.ProfileIMGURL,
			ScreenName:    row.ScreenName,
		}

		subscription.UserList = append(subscription.UserList, u)
	}

	return subscription, t.getError()
}

func (d *UserDatastore) UpdateSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	fromDb, err := d.GetSubscription(ctx, subscription.ID)
	if err != nil {
		return subscription, err
	}

	if subscription.Equal(fromDb) {
		return subscription, t.getError()
	}

	tx := t.tx
	_, err = tx.NamedExec("UPDATE subscription SET title=:title, email=:email, day=:day WHERE id = :id", subscription)
	if err != nil {
		return subscription, t.getError()
	}

	toInsert := subscription.UserList.Diff(fromDb.UserList)
	err = insertUserList(tx, toInsert, subscription.ID)
	if err != nil {
		return subscription, t.getError()
	}

	toRemove := fromDb.UserList.Diff(subscription.UserList)
	for _, u := range toRemove {
		su, err := getSubscriptionUser(tx, u.TwitterID)
		if err != nil {
			return subscription, t.getError()
		}

		_, err = tx.Exec("DELETE FROM subscription_user_m2m WHERE subscription_id=$1 AND user_id=$2", subscription.ID, su.ID)
		if err != nil {
			return subscription, t.getError()
		}
	}
	return subscription, t.getError()
}
