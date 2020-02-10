// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import models "github.com/dmtr/mail_me_all/backend/models"
import uuid "github.com/google/uuid"

// UserDatastore is an autogenerated mock type for the UserDatastore type
type UserDatastore struct {
	mock.Mock
}

// AcquireLock provides a mock function with given fields: ctx, key
func (_m *UserDatastore) AcquireLock(ctx context.Context, key uint) (bool, error) {
	ret := _m.Called(ctx, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSubscription provides a mock function with given fields: ctx, subscription
func (_m *UserDatastore) DeleteSubscription(ctx context.Context, subscription models.Subscription) error {
	ret := _m.Called(ctx, subscription)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Subscription) error); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNewSubscriptionsUsers provides a mock function with given fields: ctx, subscriptionIDs
func (_m *UserDatastore) GetNewSubscriptionsUsers(ctx context.Context, subscriptionIDs ...uuid.UUID) (map[uuid.UUID][]string, error) {
	_va := make([]interface{}, len(subscriptionIDs))
	for _i := range subscriptionIDs {
		_va[_i] = subscriptionIDs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 map[uuid.UUID][]string
	if rf, ok := ret.Get(0).(func(context.Context, ...uuid.UUID) map[uuid.UUID][]string); ok {
		r0 = rf(ctx, subscriptionIDs...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uuid.UUID][]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...uuid.UUID) error); ok {
		r1 = rf(ctx, subscriptionIDs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReadySubscriptionsStates provides a mock function with given fields: ctx, subscriptionIDs
func (_m *UserDatastore) GetReadySubscriptionsStates(ctx context.Context, subscriptionIDs ...uuid.UUID) ([]models.SubscriptionState, error) {
	_va := make([]interface{}, len(subscriptionIDs))
	for _i := range subscriptionIDs {
		_va[_i] = subscriptionIDs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []models.SubscriptionState
	if rf, ok := ret.Get(0).(func(context.Context, ...uuid.UUID) []models.SubscriptionState); ok {
		r0 = rf(ctx, subscriptionIDs...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.SubscriptionState)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...uuid.UUID) error); ok {
		r1 = rf(ctx, subscriptionIDs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscription provides a mock function with given fields: ctx, subscriptionID
func (_m *UserDatastore) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (models.Subscription, error) {
	ret := _m.Called(ctx, subscriptionID)

	var r0 models.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Subscription); ok {
		r0 = rf(ctx, subscriptionID)
	} else {
		r0 = ret.Get(0).(models.Subscription)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionTweets provides a mock function with given fields: ctx, subscriptionStateID
func (_m *UserDatastore) GetSubscriptionTweets(ctx context.Context, subscriptionStateID uint) ([]models.Tweet, error) {
	ret := _m.Called(ctx, subscriptionStateID)

	var r0 []models.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, uint) []models.Tweet); ok {
		r0 = rf(ctx, subscriptionStateID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Tweet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, subscriptionStateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionUserTweets provides a mock function with given fields: ctx, subscriptionID
func (_m *UserDatastore) GetSubscriptionUserTweets(ctx context.Context, subscriptionID uuid.UUID) (models.SubscriptionUserTweets, error) {
	ret := _m.Called(ctx, subscriptionID)

	var r0 models.SubscriptionUserTweets
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.SubscriptionUserTweets); ok {
		r0 = rf(ctx, subscriptionID)
	} else {
		r0 = ret.Get(0).(models.SubscriptionUserTweets)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptions provides a mock function with given fields: ctx, userID
func (_m *UserDatastore) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	ret := _m.Called(ctx, userID)

	var r0 []models.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.Subscription); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTodaySubscriptionsIDs provides a mock function with given fields: ctx
func (_m *UserDatastore) GetTodaySubscriptionsIDs(ctx context.Context) ([]uuid.UUID, error) {
	ret := _m.Called(ctx)

	var r0 []uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context) []uuid.UUID); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTwitterUser provides a mock function with given fields: ctx, userID
func (_m *UserDatastore) GetTwitterUser(ctx context.Context, userID uuid.UUID) (models.TwitterUser, error) {
	ret := _m.Called(ctx, userID)

	var r0 models.TwitterUser
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.TwitterUser); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(models.TwitterUser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTwitterUserByID provides a mock function with given fields: ctx, twitterUserID
func (_m *UserDatastore) GetTwitterUserByID(ctx context.Context, twitterUserID string) (models.TwitterUser, error) {
	ret := _m.Called(ctx, twitterUserID)

	var r0 models.TwitterUser
	if rf, ok := ret.Get(0).(func(context.Context, string) models.TwitterUser); ok {
		r0 = rf(ctx, twitterUserID)
	} else {
		r0 = ret.Get(0).(models.TwitterUser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, twitterUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userID
func (_m *UserDatastore) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.User); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserEmail provides a mock function with given fields: ctx, userEmail
func (_m *UserDatastore) GetUserEmail(ctx context.Context, userEmail models.UserEmail) (models.UserEmail, error) {
	ret := _m.Called(ctx, userEmail)

	var r0 models.UserEmail
	if rf, ok := ret.Get(0).(func(context.Context, models.UserEmail) models.UserEmail); ok {
		r0 = rf(ctx, userEmail)
	} else {
		r0 = ret.Get(0).(models.UserEmail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.UserEmail) error); ok {
		r1 = rf(ctx, userEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserEmails provides a mock function with given fields: ctx, status
func (_m *UserDatastore) GetUserEmails(ctx context.Context, status string) ([]models.UserEmail, error) {
	ret := _m.Called(ctx, status)

	var r0 []models.UserEmail
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.UserEmail); ok {
		r0 = rf(ctx, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.UserEmail)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertSubscription provides a mock function with given fields: ctx, subscription
func (_m *UserDatastore) InsertSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	ret := _m.Called(ctx, subscription)

	var r0 models.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, models.Subscription) models.Subscription); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Get(0).(models.Subscription)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Subscription) error); ok {
		r1 = rf(ctx, subscription)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertSubscriptionState provides a mock function with given fields: ctx, state
func (_m *UserDatastore) InsertSubscriptionState(ctx context.Context, state models.SubscriptionState) (models.SubscriptionState, error) {
	ret := _m.Called(ctx, state)

	var r0 models.SubscriptionState
	if rf, ok := ret.Get(0).(func(context.Context, models.SubscriptionState) models.SubscriptionState); ok {
		r0 = rf(ctx, state)
	} else {
		r0 = ret.Get(0).(models.SubscriptionState)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SubscriptionState) error); ok {
		r1 = rf(ctx, state)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertSubscriptionUserState provides a mock function with given fields: ctx, subscriptionID, userTwitterID, lastTweetID
func (_m *UserDatastore) InsertSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID string, lastTweetID string) error {
	ret := _m.Called(ctx, subscriptionID, userTwitterID, lastTweetID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string, string) error); ok {
		r0 = rf(ctx, subscriptionID, userTwitterID, lastTweetID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertTweet provides a mock function with given fields: ctx, tweet, subscriptionStateID
func (_m *UserDatastore) InsertTweet(ctx context.Context, tweet models.Tweet, subscriptionStateID uint) (models.Tweet, error) {
	ret := _m.Called(ctx, tweet, subscriptionStateID)

	var r0 models.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, models.Tweet, uint) models.Tweet); ok {
		r0 = rf(ctx, tweet, subscriptionStateID)
	} else {
		r0 = ret.Get(0).(models.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Tweet, uint) error); ok {
		r1 = rf(ctx, tweet, subscriptionStateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTwitterUser provides a mock function with given fields: ctx, twitterUser
func (_m *UserDatastore) InsertTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	ret := _m.Called(ctx, twitterUser)

	var r0 models.TwitterUser
	if rf, ok := ret.Get(0).(func(context.Context, models.TwitterUser) models.TwitterUser); ok {
		r0 = rf(ctx, twitterUser)
	} else {
		r0 = ret.Get(0).(models.TwitterUser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.TwitterUser) error); ok {
		r1 = rf(ctx, twitterUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUser provides a mock function with given fields: ctx, user
func (_m *UserDatastore) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUserEmail provides a mock function with given fields: ctx, userEmail
func (_m *UserDatastore) InsertUserEmail(ctx context.Context, userEmail models.UserEmail) (models.UserEmail, error) {
	ret := _m.Called(ctx, userEmail)

	var r0 models.UserEmail
	if rf, ok := ret.Get(0).(func(context.Context, models.UserEmail) models.UserEmail); ok {
		r0 = rf(ctx, userEmail)
	} else {
		r0 = ret.Get(0).(models.UserEmail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.UserEmail) error); ok {
		r1 = rf(ctx, userEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReleaseLock provides a mock function with given fields: ctx, key
func (_m *UserDatastore) ReleaseLock(ctx context.Context, key uint) (bool, error) {
	ret := _m.Called(ctx, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveOldTweets provides a mock function with given fields: ctx, tweetTTL
func (_m *UserDatastore) RemoveOldTweets(ctx context.Context, tweetTTL int) error {
	ret := _m.Called(ctx, tweetTTL)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, tweetTTL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUser provides a mock function with given fields: ctx, userID
func (_m *UserDatastore) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSubscription provides a mock function with given fields: ctx, subscription
func (_m *UserDatastore) UpdateSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	ret := _m.Called(ctx, subscription)

	var r0 models.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, models.Subscription) models.Subscription); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Get(0).(models.Subscription)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Subscription) error); ok {
		r1 = rf(ctx, subscription)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSubscriptionState provides a mock function with given fields: ctx, state
func (_m *UserDatastore) UpdateSubscriptionState(ctx context.Context, state models.SubscriptionState) (models.SubscriptionState, error) {
	ret := _m.Called(ctx, state)

	var r0 models.SubscriptionState
	if rf, ok := ret.Get(0).(func(context.Context, models.SubscriptionState) models.SubscriptionState); ok {
		r0 = rf(ctx, state)
	} else {
		r0 = ret.Get(0).(models.SubscriptionState)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SubscriptionState) error); ok {
		r1 = rf(ctx, state)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSubscriptionUserState provides a mock function with given fields: ctx, subscriptionID, userTwitterID, lastTweetID
func (_m *UserDatastore) UpdateSubscriptionUserState(ctx context.Context, subscriptionID uuid.UUID, userTwitterID string, lastTweetID string) error {
	ret := _m.Called(ctx, subscriptionID, userTwitterID, lastTweetID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string, string) error); ok {
		r0 = rf(ctx, subscriptionID, userTwitterID, lastTweetID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSubscriptionUserStateTweets provides a mock function with given fields: ctx
func (_m *UserDatastore) UpdateSubscriptionUserStateTweets(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTwitterUser provides a mock function with given fields: ctx, twitterUser
func (_m *UserDatastore) UpdateTwitterUser(ctx context.Context, twitterUser models.TwitterUser) (models.TwitterUser, error) {
	ret := _m.Called(ctx, twitterUser)

	var r0 models.TwitterUser
	if rf, ok := ret.Get(0).(func(context.Context, models.TwitterUser) models.TwitterUser); ok {
		r0 = rf(ctx, twitterUser)
	} else {
		r0 = ret.Get(0).(models.TwitterUser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.TwitterUser) error); ok {
		r1 = rf(ctx, twitterUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, user
func (_m *UserDatastore) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserEmail provides a mock function with given fields: ctx, userEmail
func (_m *UserDatastore) UpdateUserEmail(ctx context.Context, userEmail models.UserEmail) (models.UserEmail, error) {
	ret := _m.Called(ctx, userEmail)

	var r0 models.UserEmail
	if rf, ok := ret.Get(0).(func(context.Context, models.UserEmail) models.UserEmail); ok {
		r0 = rf(ctx, userEmail)
	} else {
		r0 = ret.Get(0).(models.UserEmail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.UserEmail) error); ok {
		r1 = rf(ctx, userEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
