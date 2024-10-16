package api

import (
	"errors"
	"fmt"
	"hoteRes/db"
	"hoteRes/middleware"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.store.Users.GetByFilter(c.Context(), bson.M{"email": authParams.Email})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}

	tokenStr := middleware.CreateTokenFromUser(user)
	authResponse := AuthResponse{
		User:  user,
		Token: tokenStr,
	}

	return c.JSON(authResponse)
}
