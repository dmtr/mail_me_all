package errors

import "github.com/dmtr/mail_me_all/backend/db"

type ErrorCode int

const (
	UnknownError    ErrorCode = iota
	ServerError     ErrorCode = iota
	BadRequest      ErrorCode = iota
	CantGetToken    ErrorCode = iota
	CantGetUserInfo ErrorCode = iota
	DbError         ErrorCode = iota
	NotFound        ErrorCode = iota
)

func GetErrorCode(err error) ErrorCode {
	var code ErrorCode
	switch e := err.(type) {
	case *db.DbError:
		if e.HasNoRows() == true {
			code = NotFound
		} else {
			code = DbError
		}
	default:
		code = UnknownError
	}
	return code
}
