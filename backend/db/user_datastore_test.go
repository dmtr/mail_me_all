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

	res, err := d.GetUserByID(ctx, uid)
	e, _ := err.(*DbError)
	assert.True(t, e.HasNoRows())
	assert.Empty(t, res)

	user := models.User{
		Name: "Test",
		FbID: "foo-id",
	}
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.FbID, u.FbID)

	res, err = d.GetUserByID(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u, res)

	res, err = d.GetUserByFbID(ctx, u.FbID)
	assert.NoError(t, err)
	assert.Equal(t, u, res)
}

func testUpdateUser(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	ctx := context.WithValue(context.Background(), "Tx", tx)
	user := models.User{
		Name: "Test",
		FbID: "foo-id",
	}
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, u.Name)
	assert.Equal(t, user.FbID, u.FbID)

	u.Email = "test@me.com"
	res, err := d.UpdateUser(ctx, u)
	assert.NoError(t, err)
	assert.Equal(t, u, res)

	res, err = d.GetUserByFbID(ctx, u.FbID)
	assert.NoError(t, err)
	assert.Equal(t, u, res)
}

func testInsertUserToken(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	user := models.User{
		Name: "Test",
		FbID: "some-id",
	}
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)

	token := models.Token{
		UserID:    u.ID,
		FbToken:   "some-token",
		ExpiresAt: time.Now().UTC(),
	}
	res, err := d.InsertToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, token, res)

	fromDb, err := d.GetToken(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, token, fromDb)
}

func testInsertAndUpdateUserToken(t *testing.T, tx *sqlx.Tx, d *UserDatastore) {
	user := models.User{
		Name: "Test",
		FbID: "another-id",
	}
	ctx := context.WithValue(context.Background(), "Tx", tx)
	u, err := d.InsertUser(ctx, user)
	assert.NoError(t, err)

	token := models.Token{
		UserID:    u.ID,
		FbToken:   "some-token",
		ExpiresAt: time.Now().UTC(),
	}
	res, err := d.InsertToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, token, res)

	token.FbToken = "another-token"
	updated, err := d.UpdateToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, token, updated)

	fromDb, err := d.GetToken(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, token, fromDb)
}

func TestUserDatastore(t *testing.T) {
	tests := map[string]testFunc{
		"TestInsertUserToken":          testInsertUserToken,
		"TestInsertAndUpdateUserToken": testInsertAndUpdateUserToken,
		"TestGetUser":                  testGetUser,
		"TestUpdateUser":               testUpdateUser,
	}
	runTests(tests, t)
}
