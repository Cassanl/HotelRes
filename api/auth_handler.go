package api

import (
	"errors"
	"fmt"
	"hoteRes/db"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store db.UserStore
}

func NewAuthHandler(store db.UserStore) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.store.GetByFilter(c.Context(), types.Filter{"email": authParams.Email})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}

	tokenStr := CreateTokenFromUser(user)

	c.Response().Header.Set("Authorization", tokenStr)
	c.Response().SetStatusCode(fiber.StatusNoContent)

	return nil
}
