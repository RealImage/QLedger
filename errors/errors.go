package errors

import "fmt"

type ApplicationError interface {
	//TODO Should we embed log.Error in here and add a new function New(err errors, code, message string) *ApplicationError
	Error() string
	ErrorCode() string
	ErrorMessage() string
}

//TODO create a String() method for ApplicationError

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
