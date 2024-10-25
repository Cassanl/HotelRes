package api

import (
	"hoteRes/db"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	// var params types.HotelQueryParams
	// if err := c.QueryParser(&params); err != nil {
	// 	return err
	// }

	hotels, err := h.store.Hotels.ListHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	rooms, err := h.store.Rooms.ListByFilter(c.Context(), types.Map{"hotelID": oid})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

// TODO build filter with query params
func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hotel, err := h.store.Hotels.GetByFilter(c.Context(), types.Map{"_id": oid})
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}
