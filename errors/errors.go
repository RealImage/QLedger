package errors

import "fmt"

type ApplicationError interface {
	Error() string
	ErrorCode() string
	ErrorMessage() string
}
type BaseApplicationError struct {
	Message string
	Code    string
}

func (e *BaseApplicationError) ErrorCode() string {
	return e.Code
}
func (e *BaseApplicationError) ErrorMessage() string {
	return e.Message
}
func (e *BaseApplicationError) Error() string {
	return fmt.Sprintf("%v (%v)", e.Message, e.Code)
}
