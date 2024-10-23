package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Err)
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized",
	}
}
