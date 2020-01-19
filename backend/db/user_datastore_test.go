package db

import (
	"context"
	"testing"
	"time"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
)

const (
	retry time.Duration = 4
)

type testFunc func(t *testing.T, tx *sqlx.Tx, d *UserDatastore)

func runTests(tests map[string]testFunc, t *testing.T) {

	conf := config.GetConfig()
	db, err := ConnectDb(conf.DSN, retry*time.Second)
	if err != nil {
		t.Fatal("Can't connect to database")
	}
	defer db.Close()

	d := NewUserDatastore(db)

	for name, fn := range tests {
		tx := db.MustBegin()
		f := func(t *testing.T) {
			fn(t, tx, d)
		}
		t.Run(name, f)
		tx.Rollback()
	}
}

func insertUser(d *UserDatastore, ctx context.Context) (models.User, error) {
	user := models.User{
		Name:  "Test",
		Email: "foo@bar.com",
	}
	return d.InsertUser(ctx, user)
}

func insertUserAndSubscription(d *UserDatastore, ctx context.Context) (models.User, models.Subscription, error) {
	u, err := insertUser(d, ctx)
	if err != nil {
		return models.User{}, models.Subscription{}, err
	}
	s := models.Subscription{
		UserID: u.ID,
		Title:  "test",
		Email:  "test@mail.com",
		Day:    "monday",
		UserList: []models.TwitterUserSearchResult{
			models.TwitterUserSearchResult{TwitterID: "121", Name: "foo", ProfileIMGURL: "some_url", ScreenName: "foo_name"},
			models.TwitterUserSearchResult{TwitterID: "322", Name: "bar", ProfileIMGURL: "other_url", ScreenName: "bar_name"}},
	}

	fromDb, err := d.InsertSubscription(ctx, s)
	return u, fromDb, err
}

func testGetUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	uid := uuid.New()
	ctx := context.WithValue(context.Background(), "Tx", tx)

	res, err := d.GetUser(ctx, uid)
	e, _ := err.(*DbError)
	assert.True(t, e.HasNoRows())
	assert.Empty(t, res)

	user := models.User{
		Name:  "Test",
		Email: "foo@bar.com",
	}
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.Email, u.Email)

	res, err = d.GetUser(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u, res)
}

func testUpdateUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	user := models.User{
		Name:  "Test",
		Email: "foo@bar.com",
	}
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.Email, u.Email)

	u.Email = "test@me.com"
	res, err := d.UpdateUser(ctx, u)
	assert.NoError(t, err)
	assert.Equal(t, u, res)

	res, err = d.GetUser(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u, res)
}

func testInsertTwitterUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	twitterUser := models.TwitterUser{
		UserID:        u.ID,
		TwitterID:     "111",
		AccessToken:   "some-token",
		TokenSecret:   "some-secret",
		ProfileIMGURL: "https://some_url",
	}
	res, err := d.InsertTwitterUser(ctx, twitterUser)
	assert.NoError(t, err)
	assert.Equal(t, twitterUser, res)

	fromDb, err := d.GetTwitterUserByID(ctx, twitterUser.TwitterID)
	assert.NoError(t, err)
	assert.Equal(t, twitterUser, fromDb)
}

func testInsertAndUpdateTwitterUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	twitterUser := models.TwitterUser{
		UserID:      u.ID,
		TwitterID:   "111",
		AccessToken: "some-token",
		TokenSecret: "some-secret",
	}

	res, err := d.InsertTwitterUser(ctx, twitterUser)
	assert.NoError(t, err)
	assert.Equal(t, twitterUser, res)

	twitterUser.AccessToken = "new-token"
	twitterUser.ProfileIMGURL = "new-url"
	updated, err := d.UpdateTwitterUser(ctx, twitterUser)
	assert.NoError(t, err)
	assert.Equal(t, twitterUser, updated)

	fromDb, err := d.GetTwitterUser(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, twitterUser, fromDb)
}

func testRemoveUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	twitterUser := models.TwitterUser{
		UserID:        u.ID,
		TwitterID:     "111",
		AccessToken:   "some-token",
		TokenSecret:   "some-secret",
		ProfileIMGURL: "https://some_url",
	}
	_, err = d.InsertTwitterUser(ctx, twitterUser)
	assert.NoError(t, err)

	err = d.RemoveUser(ctx, u.ID)
	assert.NoError(t, err)

	_, err = d.GetUser(ctx, u.ID)
	assert.Error(t, err)
	e, _ := err.(*DbError)
	assert.True(t, e.HasNoRows())

	_, err = d.GetTwitterUserByID(ctx, twitterUser.TwitterID)
	assert.Error(t, err)
	e, _ = err.(*DbError)
	assert.True(t, e.HasNoRows())
}

func testInsertSubscription(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	s := models.Subscription{
		UserID: u.ID,
		Title:  "test",
		Email:  "test@mail.com",
		Day:    "monday",
		UserList: []models.TwitterUserSearchResult{
			models.TwitterUserSearchResult{TwitterID: "121", Name: "foo", ProfileIMGURL: "some_url", ScreenName: "foo_name"},
			models.TwitterUserSearchResult{TwitterID: "322", Name: "bar", ProfileIMGURL: "other_url", ScreenName: "bar_name"}},
	}

	res, err := d.InsertSubscription(ctx, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.ID)
	assert.Equal(t, s.UserID, res.UserID)
	assert.Equal(t, s.Title, res.Title)
	assert.Equal(t, s.Email, res.Email)
	assert.Equal(t, s.Day, res.Day)
	assert.Equal(t, s.UserList, res.UserList)

	saved, err := d.GetSubscription(ctx, res.ID)
	assert.NoError(t, err)
	assert.Equal(t, res, saved)

	_, err = d.InsertSubscription(ctx, s)
	assert.NoError(t, err)

	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user_m2m")
	assert.NoError(t, err)
	assert.Equal(t, 4, count)

	fromDb, err := d.GetSubscriptions(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fromDb))
	assert.Equal(t, res, fromDb[0])
}

func testUpdateSubscription(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	s := models.Subscription{
		UserID: u.ID,
		Title:  "test",
		Email:  "test@mail.com",
		Day:    "monday",
		UserList: []models.TwitterUserSearchResult{
			models.TwitterUserSearchResult{TwitterID: "121", Name: "foo", ProfileIMGURL: "some_url", ScreenName: "foo_name"},
			models.TwitterUserSearchResult{TwitterID: "322", Name: "bar", ProfileIMGURL: "other_url", ScreenName: "bar_name"}},
	}

	fromDb, err := d.InsertSubscription(ctx, s)
	assert.NoError(t, err)

	s.ID = fromDb.ID
	s.Title = "test2"
	fromDb, err = d.UpdateSubscription(ctx, s)
	assert.NoError(t, err)
	assert.True(t, s.Equal(fromDb))

	s.UserList = s.UserList[1:]
	fromDb, err = d.UpdateSubscription(ctx, s)
	assert.NoError(t, err)
	assert.True(t, s.Equal(fromDb))
	assert.Equal(t, 1, len(fromDb.UserList))

	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user_m2m")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	fromDb, err = d.GetSubscription(ctx, s.ID)
	assert.NoError(t, err)
	assert.Equal(t, "322", fromDb.UserList[0].TwitterID)
	assert.Equal(t, s, fromDb)
}

func testDeleteSubscription(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	_, subscription, err := insertUserAndSubscription(d, ctx)
	assert.NoError(t, err)

	err = d.DeleteSubscription(ctx, subscription)
	assert.NoError(t, err)

	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user_m2m")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	err = tx.Get(&count, "SELECT COUNT(*) FROM subscription_user")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func testGetNewSubscriptionsUsers(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	res, err := d.GetNewSubscriptionsUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 0)

	_, s, err := insertUserAndSubscription(d, ctx)

	res, err = d.GetNewSubscriptionsUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 1)
	users, ok := res[s.ID]
	assert.True(t, ok)
	assert.Equal(t, 2, len(users))
	assert.Contains(t, users, s.UserList[0].TwitterID)
	assert.Contains(t, users, s.UserList[1].TwitterID)

	_, s2, err := insertUserAndSubscription(d, ctx)
	assert.NoError(t, err)

	res, err = d.GetNewSubscriptionsUsers(ctx, s2.ID)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 1)
	_, ok = res[s2.ID]
	assert.True(t, ok)

	res, err = d.GetNewSubscriptionsUsers(ctx, s2.ID, s.ID)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 2)
	_, ok = res[s.ID]
	assert.True(t, ok)
}

func testInsertSubscriptionState(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	_, s, err := insertUserAndSubscription(d, ctx)
	assert.NoError(t, err)

	state, err := d.InsertSubscriptionState(ctx, models.SubscriptionState{SubscriptionID: s.ID, Status: "PREPARING"})
	assert.NoError(t, err)
	assert.NotEmpty(t, state.ID)
	assert.NotEmpty(t, state.CreatedAt)
}

func testGetSubscriptionUserTweets(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	_, s, err := insertUserAndSubscription(d, ctx)
	assert.NoError(t, err)

	err = d.InsertSubscriptionUserState(ctx, s.ID, s.UserList[0].TwitterID, "123")
	assert.NoError(t, err)

	stweets, err := d.GetSubscriptionUserTweets(ctx, s.ID)
	assert.NoError(t, err)
	assert.Equal(t, s.ID, stweets.SubscriptionID)
	assert.Len(t, stweets.Tweets, 1)
	tw, ok := stweets.Tweets[s.UserList[0].TwitterID]
	assert.True(t, ok)
	assert.Equal(t, s.UserList[0].ScreenName, tw.ScreenName)
	assert.Equal(t, "123", tw.LastTweetID)
}

func testInsertUserEmail(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := insertUser(d, ctx)
	assert.NoError(t, err)

	userEmail := models.UserEmail{
		UserID: u.ID,
		Email:  "test@example.com",
		Status: models.EmailStatusNew,
	}
	res, err := d.InsertUserEmail(ctx, userEmail)
	assert.NoError(t, err)
	assert.Equal(t, userEmail, res)

	fromDb, err := d.GetUserEmail(ctx, userEmail)
	assert.NoError(t, err)
	assert.Equal(t, userEmail, fromDb)

	userEmail.Status = models.EmailStatusConfirmed
	res, err = d.UpdateUserEmail(ctx, userEmail)
	assert.NoError(t, err)
	assert.Equal(t, userEmail, res)

	emails, err := d.GetUserEmails(ctx, models.EmailStatusSent)
	assert.NoError(t, err)
	assert.Len(t, emails, 0)

	emails, err = d.GetUserEmails(ctx, models.EmailStatusConfirmed)
	assert.NoError(t, err)
	assert.Len(t, emails, 1)
}

func TestUserDatastore(t *testing.T) {
	tests := map[string]testFunc{
		"TestInsertTwitterUser":          testInsertTwitterUser,
		"TestInsertAndUpdateTwitterUser": testInsertAndUpdateTwitterUser,
		"TestGetUser":                    testGetUser,
		"TestUpdateUser":                 testUpdateUser,
		"TestRemoveUser":                 testRemoveUser,
		"TestInsertSubscription":         testInsertSubscription,
		"TestUpdatetSubscription":        testUpdateSubscription,
		"TestDeleteSubscription":         testDeleteSubscription,
		"TestGetNewSubscriptionsUsers":   testGetNewSubscriptionsUsers,
		"TestInsertSubscriptionState":    testInsertSubscriptionState,
		"TestGetSubscriptionUserTweets":  testGetSubscriptionUserTweets,
		"TestInsertUserEmail":            testInsertUserEmail,
	}
	runTests(tests, t)
}
