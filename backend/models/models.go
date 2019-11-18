package models

import (
	"context"
	"fmt"
	"sort"

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
	ScreenName    string
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

// TwitterUser - represents twitter account
type TwitterUserSearchResult struct {
	Name          string `db:"name"`
	TwitterID     string `db:"twitter_id"`
	ProfileIMGURL string `db:"profile_image_url"`
	ScreenName    string `db:"screen_name"`
}

func (t TwitterUserSearchResult) String() string {
	return fmt.Sprintf("TwitterUserSearchResult: TwitterID %s, Name %s", t.TwitterID, t.Name)
}

func (t TwitterUserSearchResult) Equal(another TwitterUserSearchResult) bool {
	return t.TwitterID == another.TwitterID
}

type UserList []TwitterUserSearchResult

func (u UserList) Len() int           { return len(u) }
func (u UserList) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u UserList) Less(i, j int) bool { return u[i].TwitterID < u[j].TwitterID }

func (u UserList) Diff(another UserList) UserList {
	sorted := append(another[:0:0], another...)
	sort.Sort(sorted)
	res := make(UserList, 0, len(sorted))
	length := len(sorted)
	for _, user := range u {
		i := sort.Search(length, func(i int) bool { return sorted[i].TwitterID >= user.TwitterID })
		if i == length || !user.Equal(sorted[i]) {
			res = append(res, user)
		}
	}
	return res
}

// Subscription represents user subscription
type Subscription struct {
	ID       uuid.UUID `db:"id"`
	UserID   uuid.UUID `db:"user_id"`
	Title    string    `db:"title"`
	Email    string    `db:"email"`
	Day      string    `db:"day"`
	UserList UserList
}

func (s Subscription) String() string {
	return fmt.Sprintf("Subscription: ID %s, UserID %s, Title %s, users amount %d", s.ID, s.UserID, s.Title, len(s.UserList))
}

func (s Subscription) Equal(another Subscription) bool {
	if s.ID != another.ID {
		return false
	}

	if s.UserID != another.UserID {
		return false
	}

	if s.Title != another.Title {
		return false
	}

	if s.Email != another.Email {
		return false
	}

	if s.Day != another.Day {
		return false
	}

	if len(s.UserList) != len(another.UserList) {
		return false
	}

	diff := s.UserList.Diff(another.UserList)
	if len(diff) != 0 {
		return false
	}

	return true
}

// UserUseCase - represents user use cases
type UserUseCase interface {
	SignInWithTwitter(ctx context.Context, twitterID, name, email, screenName, accessToken, tokenSecret string) (User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
	SearchTwitterUsers(ctx context.Context, userID uuid.UUID, query string) ([]TwitterUserSearchResult, error)
	AddSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]Subscription, error)
	UpdateSubscription(ctx context.Context, userID uuid.UUID, subscription Subscription) (Subscription, error)
	DeleteSubscription(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) error
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

	InsertSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]Subscription, error)
	UpdateSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (Subscription, error)
	DeleteSubscription(ctx context.Context, subscription Subscription) error

	GetNewSubscriptionsIDs(ctx context.Context) ([]uuid.UUID, error)
	InsertSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID, lastTweetID string) error
}

type SystemUseCase interface {
	InitSubscriptions(ids ...uuid.UUID) error
}

// UseCases - represents all use cases
type UseCases struct {
	UserUseCase
	SystemUseCase
}

// NewUseCases - returns new UseCases struct
func NewUseCases(user UserUseCase, system SystemUseCase) *UseCases {
	return &UseCases{user, system}
}
