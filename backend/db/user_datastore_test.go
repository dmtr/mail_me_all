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

func TestUserDatastore(t *testing.T) {
	tests := map[string]testFunc{
		"TestInsertTwitterUser":          testInsertTwitterUser,
		"TestInsertAndUpdateTwitterUser": testInsertAndUpdateTwitterUser,
		"TestGetUser":                    testGetUser,
		"TestUpdateUser":                 testUpdateUser,
	}
	runTests(tests, t)
}
