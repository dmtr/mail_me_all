package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Model interface
type Model interface {
	String() string
}

// TwitterUser - represents twitter account
type TwitterUser struct {
	UserID      uuid.UUID `db:"user_id"`
	TwitterID   string    `db:"social_account_id"`
	AccessToken string    `db:"access_token"`
	TokenSecret string    `db:"token_secret"`
}

func (t TwitterUser) String() string {
	return fmt.Sprintf("TwitterUser: UserID %s, TwitterID %s", t.UserID, t.TwitterID)
}

// User - represents user
type User struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("User: Name %s, ID %s", u.Name, u.ID)
}

// UserUseCase - represents user use cases
type UserUseCase interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
}

type UserDatastore interface {
	InsertUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (User, error)

	InsertTwitterUser(ctx context.Context, twitterUser TwitterUser) (TwitterUser, error)
	UpdateTwitterUser(ctx context.Context, twitterUser TwitterUser) (TwitterUser, error)
	GetTwitterUserByID(ctx context.Context, twitterUserID string) (TwitterUser, error)
	GetTwitterUser(ctx context.Context, userID uuid.UUID) (TwitterUser, error)
}

type UseCases struct {
	User UserUseCase
}

func NewUseCases(user UserUseCase) *UseCases {
	return &UseCases{User: user}
}
