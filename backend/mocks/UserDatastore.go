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
