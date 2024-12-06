package custom_errors

import "fmt"

type InvalidFormat struct {
	msg       string
	bad_field string
}

func (e *InvalidFormat) Error() string {
	return fmt.Sprint("Invalid format: ", e.msg, " in field: ", e.bad_field)
}

func NewInvalidFormat(msg, bad_field string) *InvalidFormat {
	return &InvalidFormat{msg: msg, bad_field: bad_field}
}

type UnAuthorized struct {
	msg string
}

func (e *UnAuthorized) Error() string {
	return fmt.Sprint("UnAuthorized: ", e.msg)
}

func NewUnAuthorized(msg string) *UnAuthorized {
	return &UnAuthorized{msg: msg}
}

type InternalError struct {
	msg string
}

func (e *InternalError) Error() string {
	return fmt.Sprint("Internal Error: ", e.msg)
}

func NewInternalError(msg string) *InternalError {
	return &InternalError{msg: msg}
}

type NotFound struct {
	msg string
}

func (e *NotFound) Error() string {
	return fmt.Sprint("Not Found: ", e.msg)
}

func NewNotFound(msg string) *NotFound {
	return &NotFound{msg: msg}
}

type DuplicateInformation struct {
	msg string
}

func (e *DuplicateInformation) Error() string {
	return fmt.Sprint("Duplicate Information: ", e.msg)
}

func NewDuplicateInformation(msg string) *DuplicateInformation {
	return &DuplicateInformation{msg: msg}
}
