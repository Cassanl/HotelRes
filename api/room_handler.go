package api

import (
	"hoteRes/db"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBook(c *fiber.Ctx) error {
	roomOid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	_ = roomOid

	user := c.Context().UserValue("user").(*types.User)
	_ = user

	return nil
}
