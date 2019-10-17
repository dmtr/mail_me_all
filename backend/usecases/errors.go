package usecases

import "github.com/dmtr/mail_me_all/backend/db"

type ErrorCode int

const (
	unknownError    ErrorCode = 1 << iota
	cantGetToken              = iota
	cantGetUserInfo           = iota
	dbError                   = iota
	notFound                  = iota
)

func getErrorCode(err error) ErrorCode {
	var code ErrorCode
	switch e := err.(type) {
	case *db.DbError:
		if e.HasNoRows() == true {
			code = notFound
		} else {
			code = dbError
		}
	default:
		code = unknownError
	}
	return code
}

type UseCaseError struct {
	msg  string
	code ErrorCode
}

func NewUseCaseError(msg string, code ErrorCode) *UseCaseError {
	return &UseCaseError{msg: msg, code: code}
}

func (e *UseCaseError) Error() string {
	return e.msg
}

func (e *UseCaseError) Code() ErrorCode {
	return e.code
}
