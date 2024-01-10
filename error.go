package fswap

import (
	"errors"
	"fmt"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

func IsErrorCode(err error, code int) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}
