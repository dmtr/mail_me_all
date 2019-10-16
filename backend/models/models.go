package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Model interface
type Model interface {
	String() string
}

// User - represents user
type User struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
	FbID  string    `db:"fb_id"`
}

func (u User) String() string {
	return fmt.Sprintf("User: Name %s, FbID %s", u.Name, u.FbID)
}

// Token - represents user token
type Token struct {
	UserID    uuid.UUID `db:"user_id"`
	FbToken   string    `db:"fb_token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (t Token) String() string {
	return fmt.Sprintf("Token: UserID %s, ExpiresAt %s", t.UserID, t.ExpiresAt)
}

func (t Token) CalculateExpiresAt(expiresIn uint64) time.Time {
	now := time.Now().UTC()
	return now.Add(time.Duration(expiresIn) * time.Second)
}

type UserUseCase interface {
	SignInFB(ctx context.Context, userID string, accessToken string) (User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
}

type UserDatastore interface {
	InsertUser(ctx context.Context, user User) (User, error)
	InsertToken(ctx context.Context, token Token) (Token, error)

	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
}

type UseCases struct {
	User UserUseCase
}

func NewUseCases(user UserUseCase) *UseCases {
	return &UseCases{User: user}
}
