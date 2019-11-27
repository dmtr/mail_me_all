package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

const (
	//Preparing - subscription status
	Preparing string = "PREPARING"

	//Ready - subscription status
	Ready string = "READY"

	//Sending - subscription status
	Sending string = "SENDING"

	//Sent - subscription status
	Sent string = "SENT"

	//Failed - subscription status
	Failed string = "FAILED"
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

// TwitterUserSearchResult - twitter user info
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

func (s Subscription) GetSubject() string {
	return fmt.Sprintf("New Issue of %s", s.Title)
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

// SubscriptionState - subscription status
type SubscriptionState struct {
	ID             uint      `db:"id"`
	SubscriptionID uuid.UUID `db:"subscription_id"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (s *SubscriptionState) String() string {
	return fmt.Sprintf(
		"SubscriptionState: id %d, subscription_id %s, status %s, updated at %s", s.ID, s.SubscriptionID, s.Status, s.UpdatedAt)
}

// UserTweet - last read tweet of a user
type UserLastTweet struct {
	ScreenName  string
	LastTweetID string
}

func (u UserLastTweet) String() string {
	return fmt.Sprintf("UserLastTweet: ScreenName %s, LastTweetID %s", u.ScreenName, u.LastTweetID)
}

// SubscriptionUserTweets contains map user twiiter id to UserLastTweet
type SubscriptionUserTweets struct {
	SubscriptionID uuid.UUID `db:"subscription_id"`
	Tweets         map[string]UserLastTweet
}

func (s SubscriptionUserTweets) String() string {
	return fmt.Sprintf("SubscriptionUserTweets: SubscriptionID %s", s.SubscriptionID)
}

// Tweet - user tweet
type Tweet struct {
	ID      uint       `db:"id"`
	TweetID string     `db:"tweet_id"`
	Tweet   TweetAttrs `db:"tweet"`
}

func (t Tweet) String() string {
	return fmt.Sprintf("Tweet %s", t.TweetID)
}

//TweetAttrs - tweet data
type TweetAttrs struct {
	IdStr                string `json:"id_str"`
	Text                 string `json:"text"`
	FullText             string `json:"full_text"`
	InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
	InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
	UserId               string `json:"user_id"`
	UserName             string `json:"user_name"`
	UserScreenName       string `json:"user_screen_name"`
	UserProfileImageUrl  string `json:"user_profile_image_url"`
}

func (a TweetAttrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *TweetAttrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
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

	GetNewSubscriptionsUsers(ctx context.Context, subscriptionIDs ...uuid.UUID) (map[uuid.UUID][]string, error)
	InsertSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID, lastTweetID string) error
	GetTodaySubscriptionsIDs(ctx context.Context) ([]uuid.UUID, error)
	InsertSubscriptionState(ctx context.Context, state SubscriptionState) (SubscriptionState, error)
	UpdateSubscriptionState(ctx context.Context, state SubscriptionState) (SubscriptionState, error)
	GetReadySubscriptionsStates(ctx context.Context, subscriptionIDs ...uuid.UUID) ([]SubscriptionState, error)

	GetSubscriptionUserTweets(ctx context.Context, subscriptionID uuid.UUID) (SubscriptionUserTweets, error)
	GetSubscriptionTweets(ctx context.Context, subscriptionStateID uint) ([]Tweet, error)

	InsertTweet(ctx context.Context, tweet Tweet, subscriptionStateID uint) (Tweet, error)
}

type SystemUseCase interface {
	InitSubscriptions(ids ...uuid.UUID) error
	PrepareSubscriptions(ids ...uuid.UUID) error
	SendSubscriptions(templatePath string, ids ...uuid.UUID) error
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
