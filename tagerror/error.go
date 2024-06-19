package tagerror

import "fmt"

type TagError struct {
	Code int
	Err  error
	Msg  string
}

var ErrNoMetadata = 001

func (e *TagError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Original error: %v", e.Code, e.Msg, e.Err)
}

func (e *TagError) Is(code int) bool {
	return e.Code == code
}

func NewTagError(code int, msg string, err error) error {
	return &TagError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}
