package fswap

import (
	"fmt"
)

type Error struct {
	Code int `json:"code,omitempty"`
	Msg string `json:"msg,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s",err.Code,err.Msg)
}
