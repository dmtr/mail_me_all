package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const (
	uniqueViolationErr pq.ErrorCode = "23505"
)

// DbError contains details about db error
type DbError struct {
	PqError *pq.Error
	Err     error
}

func (e *DbError) Error() string {
	if e.PqError != nil {
		return fmt.Sprintf("Got Database error with code: %s message: %s detail: %s and constraint: %s",
			e.PqError.Code, e.PqError.Message, e.PqError.Detail, e.PqError.Constraint)
	}
	return e.Err.Error()
}

func (e *DbError) HasNoRows() bool {
	return e.Err == sql.ErrNoRows
}

func (e *DbError) IsUniqueViolationError() bool {
	if e.PqError != nil {
		return e.PqError.Code == uniqueViolationErr
	}
	return false
}

func getPqError(err error) *pq.Error {
	if err, ok := err.(*pq.Error); ok {
		return err
	} else {
		return nil
	}
}

func getDbError(err error) error {
	if err != nil {
		return &DbError{
			PqError: getPqError(err),
			Err:     err,
		}
	}
	return err
}
