package models

// User - represents user
type User struct {
	Name    string `json:"name" db:"name" binding:"required"`
	FbID    string `json:"fbid" db:"fb_id"`
	FbToken string `json:"fbtoken" db:"fb_token"`
}
