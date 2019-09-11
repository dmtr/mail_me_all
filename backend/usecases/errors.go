package usecases

type UseCaseError struct {
	msg string
}

func NewUseCaseError(msg string) *UseCaseError {
	return &UseCaseError{msg: msg}
}

func (e *UseCaseError) Error() string {
	return e.msg
}
