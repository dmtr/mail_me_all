package db

import (
	"context"
	"fmt"
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
		fmt.Printf("Running test %s", name)
		tx := db.MustBegin()
		f := func(t *testing.T) {
			fn(t, tx, d)
		}
		t.Run(name, f)
		tx.Rollback()
	}
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
	user := models.User{
		Name:  "Test",
		Email: "some@body.com",
	}
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := d.InsertUser(ctx, user)
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
	user := models.User{
		Name:  "Test",
		Email: "another@mail.com",
	}
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := d.InsertUser(ctx, user)
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

func testInsertSubscription(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	user := models.User{
		Name:  "Test",
		Email: "another@mail.com",
	}
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := d.InsertUser(ctx, user)
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

func TestUserDatastore(t *testing.T) {
	tests := map[string]testFunc{
		"TestInsertTwitterUser":          testInsertTwitterUser,
		"TestInsertAndUpdateTwitterUser": testInsertAndUpdateTwitterUser,
		"TestGetUser":                    testGetUser,
		"TestUpdateUser":                 testUpdateUser,
		"TestInsertSubscription":         testInsertSubscription,
	}
	runTests(tests, t)
}
