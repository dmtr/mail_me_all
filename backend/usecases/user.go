package usecases

import "github.com/dmtr/mail_me_all/backend/models"
import log "github.com/sirupsen/logrus"

type UserUseCase struct {
	models.UserDatastore
}

func NewUserUseCase(datastore models.UserDatastore) UserUseCase {
	return UserUseCase{datastore}
}

func (u UserUseCase) CreateUser(user *models.User) error {
	log.Debugf("Going to create user %v", user)
	return u.UserDatastore.CreateUser(user)
}
