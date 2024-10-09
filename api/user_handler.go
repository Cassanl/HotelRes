package api

import (
	"context"
	"fmt"
	"hoteRes/db"
	"hoteRes/types"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("pol")
}

func (h *UserHandler) HandleInsertUser(c *fiber.Ctx) error {
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	return nil
}
