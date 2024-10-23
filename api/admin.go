package api

import (
	"fmt"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue(types.UserKey).(*types.User)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("unauthorized")
	}
	return c.Next()
}
