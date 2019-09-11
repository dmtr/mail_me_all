package models

// User - represents user
type User struct {
	Name    string `json:"name" db:"name" binding:"required"`
	FbID    string `json:"fbid" db:"fb_id"`
	FbToken string `json:"fbtoken" db:"fb_token"`
}

type UserUseCase interface {
	CreateUser(user *User) error
}

type UserDatastore interface {
	CreateUser(user *User) error
}

type UseCases struct {
	User UserUseCase
}

func NewUseCases(user UserUseCase) *UseCases {
	return &UseCases{User: user}
}
