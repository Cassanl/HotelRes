package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	if err, ok := err.(Error); ok {
		return c.Status(err.Code).JSON(err)
	}
	apiErr := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiErr.Code).JSON(apiErr)
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

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid request",
	}
}

func ErrResourceNotFound() Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  "resource not found",
	}
}
