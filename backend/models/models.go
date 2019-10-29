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
	UserID        uuid.UUID `db:"user_id"`
	TwitterID     string    `db:"social_account_id"`
	AccessToken   string    `db:"access_token"`
	TokenSecret   string    `db:"token_secret"`
	ProfileIMGURL string    `db:"profile_image_url"`
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
	SignInWithTwitter(ctx context.Context, twitterID, name, email, screenName, accessToken, tokenSecret string) (User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
}

// UserDatastore - represents all user related database methods
type UserDatastore interface {
	InsertUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (User, error)

	InsertTwitterUser(ctx context.Context, twitterUser TwitterUser) (TwitterUser, error)
	UpdateTwitterUser(ctx context.Context, twitterUser TwitterUser) (TwitterUser, error)
	GetTwitterUserByID(ctx context.Context, twitterUserID string) (TwitterUser, error)
	GetTwitterUser(ctx context.Context, userID uuid.UUID) (TwitterUser, error)
}

// UseCases - represents all use cases
type UseCases struct {
	User UserUseCase
}

// NewUseCases - returns new UseCases struct
func NewUseCases(user UserUseCase) *UseCases {
	return &UseCases{User: user}
}
