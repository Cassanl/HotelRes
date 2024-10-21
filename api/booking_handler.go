package api

import (
	"hoteRes/db"
	"hoteRes/middleware"
	"hoteRes/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetCurrentUserBookings(c *fiber.Ctx) error {
	user, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	bookings, err := h.store.Bookings.ListByFilter(c.Context(), types.Filter{"userID": user.ID})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	booking, err := h.store.Bookings.GetByFilter(c.Context(), types.Filter{"_id": oid})
	if err != nil {
		return err
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Bookings.ListByFilter(c.Context(), types.Filter{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}
