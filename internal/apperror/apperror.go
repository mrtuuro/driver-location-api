package apperror

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func NewAppError(code string, err error, publicMsg ...string) *AppError {
	msg := ""
	if len(publicMsg) > 0 {
		msg = publicMsg[0]
	}
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

func (e *AppError) Unwrap() error {
	return e.Err
}
