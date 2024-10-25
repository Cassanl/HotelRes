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

func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	room, err := h.store.Rooms.GetByFilter(c.Context(), types.Map{"_id": oid})
	if err != nil {
		return err
	}
	return c.JSON(room)
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Rooms.ListByFilter(c.Context(), types.Map{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
