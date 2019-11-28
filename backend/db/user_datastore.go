package db

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (d *UserDatastore) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	log.Debugf("Going to remove user %s", userID)
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.Exec("DELETE FROM user_account WHERE id = $1", userID)
	if err != nil {
		log.Error(err.Error() + fmt.Sprintf(" removing user: %s", userID))
		return t.getError()
	}
	return t.getError()

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

	fromDb, err := d.GetSubscription(context.WithValue(ctx, "Tx", t.tx), subscription.ID)
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

		_, err = tx.Exec("DELETE FROM subscription_user_state WHERE subscription_id=$1 AND user_twitter_id=$2", subscription.ID, su.TwitterID)
		if err != nil {
			return subscription, t.getError()
		}

	}
	return subscription, t.getError()
}

func (d *UserDatastore) DeleteSubscription(ctx context.Context, subscription models.Subscription) error {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.NamedExec("DELETE FROM subscription WHERE id = :id", subscription)
	return t.getError()
}

func (d *UserDatastore) GetNewSubscriptionsUsers(ctx context.Context, subscriptionIDs ...uuid.UUID) (map[uuid.UUID][]string, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	res := make(map[uuid.UUID][]string, 0)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.Select("s.id AS subscription_id, array_agg(t.twitter_id) AS users").
		Prefix("WITH t AS " +
			"(SELECT u.id, u.twitter_id, st.last_tweet_id " +
			"FROM subscription_user u " +
			"LEFT JOIN subscription_user_state st ON st.user_twitter_id = u.twitter_id)").
		From("t").
		JoinClause("INNER JOIN subscription_user_m2m m ON m.user_id = t.id").
		JoinClause("INNER JOIN subscription s ON s.id = m.subscription_id").
		Where("t.last_tweet_id IS NULL")

	if len(subscriptionIDs) > 0 {
		q = q.Where(sq.Eq{"subscription_id": subscriptionIDs})
	}
	q = q.GroupBy("s.id")

	sql, args, err := q.ToSql()
	if err != nil {
		return res, t.getError()
	}

	var rows *sqlx.Rows
	if len(args) > 0 {
		rows, err = t.tx.Queryx(sql, args...)
	} else {
		rows, err = t.tx.Queryx(sql)
	}

	if err != nil {
		log.Errorf("Got error %s %s", err, sql)
		return res, t.getError()
	}

	for rows.Next() {
		var subscriptionID uuid.UUID
		var users []string

		err = rows.Scan(&subscriptionID, pq.Array(&users))
		if err != nil {
			log.Errorf("Got error %s", err)
			continue
		}
		res[subscriptionID] = users
	}

	return res, t.getError()
}

func (d *UserDatastore) InsertSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID string, lastTweetID string) error {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.Exec("INSERT INTO subscription_user_state (subscription_id, user_twitter_id, last_tweet_id) VALUES ($1, $2, $3)", subscriptionID, userTwitterID, lastTweetID)

	return err
}

func (d *UserDatastore) UpdateSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID string, lastTweetID string) error {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.Exec("UPDATE subscription_user_state SET last_tweet_id = $1 WHERE subscription_id = $2 AND user_twitter_id = $3", lastTweetID, subscriptionID, userTwitterID)

	return err
}

func (d *UserDatastore) UpdateSubscriptionUserStateTweets(ctx context.Context) error {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	rows, err := t.tx.Queryx("SELECT t.tweet->>'user_id' AS user_id, st.subscription_id, MAX(t.tweet_id::BIGINT) AS tweet_id " +
		"FROM subscription_state st " +
		"INNER JOIN subscription_state_tweet_m2m m ON m.subscription_state_id = st.id " +
		"INNER JOIN tweet t ON t.id = m.tweet_id " +
		"WHERE st.status = 'SENT' AND st.created_at::DATE = NOW()::DATE " +
		"GROUP  BY t.tweet->>'user_id', subscription_id")

	if err != nil {
		return err
	}

	type row struct {
		UserID         string    `db:"user_id"`
		SubscriptionID uuid.UUID `db:"subscription_id"`
		TweetID        uint      `db:"tweet_id"`
	}

	var results []row

	for rows.Next() {
		var r row
		err = rows.StructScan(&r)
		if err != nil {
			log.Errorf("Can't scan %s", err)
			continue
		}

		results = append(results, r)
	}

	for _, r := range results {
		lastTweetID := strconv.FormatUint(uint64(r.TweetID), 10)
		err = d.UpdateSubscriptionUserState(context.WithValue(ctx, "Tx", t.tx), r.SubscriptionID, r.UserID, lastTweetID)
	}

	return err
}

func (d *UserDatastore) GetTodaySubscriptionsIDs(ctx context.Context) ([]uuid.UUID, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	rows, err := t.tx.Queryx("WITH t AS " +
		"(SELECT subscription_id FROM subscription_state st " +
		"WHERE st.created_at::DATE = NOW()::DATE) " +
		"SELECT s.id FROM subscription s " +
		"LEFT JOIN t ON s.id = t.subscription_id " +
		"WHERE s.day = get_day_of_week(NOW()) " +
		"GROUP BY s.id HAVING count(t.*) = 0",
	)

	if err != nil {
		log.Errorf("Got error %s", err)
		return []uuid.UUID{}, t.getError()
	}

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			log.Errorf("Got error %s", err)
			continue
		}
		ids = append(ids, id)
	}

	return ids, t.getError()
}

func (d *UserDatastore) UpdateSubscriptionState(ctx context.Context, state models.SubscriptionState) (models.SubscriptionState, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	_, err = t.tx.NamedExec(
		"UPDATE subscription_state SET status = (:status) WHERE id = :id ", state)
	if err != nil {
		return state, err
	}

	return state, err
}

func (d *UserDatastore) InsertSubscriptionState(ctx context.Context, state models.SubscriptionState) (models.SubscriptionState, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	res, err := t.tx.NamedQuery(
		"INSERT INTO subscription_state (subscription_id, status) VALUES (:subscription_id, :status) RETURNING id", state)
	if err != nil {
		return state, err
	}

	var id uint
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			log.Errorf("Scan error: %s", err)
			return state, err
		}
	}

	var fromDB models.SubscriptionState
	err = t.tx.Get(&fromDB, "SELECT id, subscription_id, status, created_at, updated_at FROM subscription_state WHERE id=$1", id)
	if err != nil {
		return state, err
	}

	return fromDB, err
}

func (d *UserDatastore) GetSubscriptionUserTweets(ctx context.Context, subscriptionID uuid.UUID) (models.SubscriptionUserTweets, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	rows, err := t.tx.Queryx("SELECT u.twitter_id, u.screen_name, st.last_tweet_id "+
		"FROM subscription_user u INNER JOIN subscription_user_state st ON st.user_twitter_id = u.twitter_id "+
		"WHERE st.subscription_id = $1", subscriptionID)
	if err != nil {
		return models.SubscriptionUserTweets{}, err
	}

	res := models.SubscriptionUserTweets{
		SubscriptionID: subscriptionID,
		Tweets:         make(map[string]models.UserLastTweet, 0)}

	type row struct {
		TwitterId   string `db:"twitter_id"`
		ScreenName  string `db:"screen_name"`
		LastTweetID string `db:"last_tweet_id"`
	}

	for rows.Next() {
		var r row
		err = rows.StructScan(&r)
		if err != nil {
			log.Errorf("Got error %s", err)
			continue
		}

		res.Tweets[r.TwitterId] = models.UserLastTweet{ScreenName: r.ScreenName, LastTweetID: r.LastTweetID}
	}

	return res, err
}

func (d *UserDatastore) InsertTweet(ctx context.Context, tweet models.Tweet, subscriptionStateID uint) (models.Tweet, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	tx := t.tx

	_, err = tx.Exec("SAVEPOINT save1")
	if err != nil {
		return tweet, err
	}

	res, err := tx.NamedQuery(
		"INSERT INTO tweet (tweet_id, tweet) VALUES (:tweet_id, :tweet) RETURNING id", tweet)

	if err != nil {
		e := getDbError(err).(*DbError)
		if e.IsUniqueViolationError() {
			_, err = tx.Exec("ROLLBACK TO SAVEPOINT save1")
			if err != nil {
				return tweet, err
			}

			var fromDB models.Tweet
			err = t.tx.Get(&fromDB, "SELECT id, tweet_id, tweet FROM tweet WHERE tweet_id=$1", tweet.TweetID)
			if err != nil {
				return tweet, err
			}

			tweet = fromDB
			err = nil
		} else {
			return tweet, err
		}
	}

	if tweet.ID == 0 {
		var id uint
		for res.Next() {
			err = res.Scan(&id)
			if err != nil {
				log.Errorf("Scan error: %s", err)
				return tweet, err
			}
		}

		tweet.ID = id
	}

	m2m := struct {
		SubscriptionStateID uint `db:"subscription_state_id"`
		TweetID             uint `db:"tweet_id"`
	}{
		subscriptionStateID,
		tweet.ID,
	}
	_, err = tx.NamedExec("INSERT INTO subscription_state_tweet_m2m (subscription_state_id, tweet_id) VALUES(:subscription_state_id, :tweet_id)", m2m)
	if err != nil {
		return tweet, err
	}

	return tweet, err
}

func (d *UserDatastore) GetReadySubscriptionsStates(ctx context.Context, subscriptionIDs ...uuid.UUID) ([]models.SubscriptionState, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	res := make([]models.SubscriptionState, 0)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.Select("id, subscription_id, status, created_at, updated_at FROM subscription_state").
		Where("status = 'READY'")

	if len(subscriptionIDs) > 0 {
		q = q.Where(sq.Eq{"subscription_id": subscriptionIDs})
	} else {
		q = q.Where("created_at::DATE = NOW()::DATE")
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return res, t.getError()
	}

	var rows *sqlx.Rows
	if len(args) > 0 {
		rows, err = t.tx.Queryx(sql, args...)
	} else {
		rows, err = t.tx.Queryx(sql)
	}

	if err != nil {
		log.Errorf("Got error %s, sql %s", err, sql)
		return []models.SubscriptionState{}, t.getError()
	}

	for rows.Next() {
		var s models.SubscriptionState
		err = rows.StructScan(&s)
		if err != nil {
			log.Errorf("Got error %s", err)
			continue
		}
		res = append(res, s)
	}

	return res, t.getError()
}

func (d *UserDatastore) GetSubscriptionTweets(ctx context.Context, subscriptionStateID uint) ([]models.Tweet, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	rows, err := t.tx.Queryx("SELECT t.id, t.tweet_id, t.tweet FROM tweet t "+
		"INNER JOIN subscription_state_tweet_m2m m ON t.id = m.tweet_id "+
		"WHERE m.subscription_state_id = $1 "+
		"ORDER BY t.tweet_id::BIGINT ASC", subscriptionStateID)

	if err != nil {
		log.Errorf("Got error %s", err)
		return []models.Tweet{}, t.getError()
	}

	tweets := make([]models.Tweet, 0)
	for rows.Next() {
		var t models.Tweet
		err = rows.StructScan(&t)
		if err != nil {
			log.Errorf("Got error %s", err)
			continue
		}
		tweets = append(tweets, t)
	}

	return tweets, t.getError()
}

func (d *UserDatastore) AcquireLock(ctx context.Context, key uint) (bool, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var res bool
	err = t.tx.Get(&res, "SELECT pg_try_advisory_lock($1)", key)

	return res, err
}

func (d *UserDatastore) ReleaseLock(ctx context.Context, key uint) (bool, error) {
	var err error
	t := getTransaction(ctx, d.DB, &err)

	defer func() {
		t.commitOrRollback()
	}()

	var res bool
	err = t.tx.Get(&res, "SELECT pg_advisory_unlock($1)", key)

	return res, err
}
