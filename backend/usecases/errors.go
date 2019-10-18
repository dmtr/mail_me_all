package usecases

import "github.com/dmtr/mail_me_all/backend/errors"

type UseCaseError struct {
	msg  string
	code errors.ErrorCode
}

func NewUseCaseError(msg string, code errors.ErrorCode) *UseCaseError {
	return &UseCaseError{msg: msg, code: code}
}

func (e *UseCaseError) Error() string {
	return e.msg
}

func (e *UseCaseError) Code() errors.ErrorCode {
	return e.code
}
