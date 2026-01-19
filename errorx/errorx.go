package errorx

import "fmt"

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCodeError(code int, message string) *CodeError {
	return &CodeError{
		Code:    code,
		Message: message,
	}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
