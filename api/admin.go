package api

import (
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue(types.UserKey).(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if !user.IsAdmin {
		return ErrUnauthorized()
	}
	return c.Next()
}
