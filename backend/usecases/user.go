package usecases

import "github.com/dmtr/mail_me_all/backend/models"
import log "github.com/sirupsen/logrus"

const (
	userCreationError string = "Can not create user"
)

type UserUseCase struct {
	models.UserDatastore
}

func NewUserUseCase(datastore models.UserDatastore) *UserUseCase {
	return &UserUseCase{datastore}
}

func (u UserUseCase) CreateUser(user *models.User) error {
	log.Debugf("Going to create user %v", user)
	if err := u.UserDatastore.CreateUser(user); err != nil {
		return NewUseCaseError(userCreationError)
	}
	return nil
}
