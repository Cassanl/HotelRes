package api

import (
	"context"
	"fmt"
	"hoteRes/db"
	"hoteRes/types"
	"net/http"
	"time"

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

func (h *BookingHandler) HandlePostBooking(c *fiber.Ctx) error {
	var params types.PostBookingParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errs := params.Validate(); len(errs) != 0 {
		return c.JSON(errs)
	}

	user, ok := c.Context().UserValue(types.UserKey).(*types.User)
	if !ok {
		return ErrUnauthorized()
	}

	booked, err := h.isBooked(c.Context(), params)
	if err != nil {
		return err
	}
	if booked {
		return NewError(http.StatusBadRequest, "room already booked")
	}

	booking := types.NewBookingFromParams(params, user.ID)
	booking, err = h.store.Bookings.Insert(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	filters := types.Map{
		"cancelled":   true,
		"cancelledAt": time.Now(),
	}
	if err := h.store.Bookings.Update(c.Context(), types.Map{"_id": oid}, filters); err != nil {
		return err
	}

	resp := types.GenericResponse{
		Status: http.StatusOK,
		Msg:    "booking cancelled",
	}
	return c.Status(resp.Status).JSON(resp.Msg)
}

func (h *BookingHandler) HandleDeleteBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.store.Bookings.Delete(c.Context(), id); err != nil {
		return err
	}

	resp := types.GenericResponse{
		Status: http.StatusOK,
		Msg:    fmt.Sprintf("deleted %s", id),
	}
	return c.Status(resp.Status).JSON(resp.Msg)
}

func (h *BookingHandler) HandleGetCurrentUserBookings(c *fiber.Ctx) error {
	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	bookings, err := h.store.Bookings.ListByFilter(c.Context(), types.Map{"userID": user.ID})
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
	booking, err := h.store.Bookings.GetByFilter(c.Context(), types.Map{"_id": oid})
	if err != nil {
		return err
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Bookings.ListByFilter(c.Context(), types.Map{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) isBooked(ctx context.Context, params types.PostBookingParams) (bool, error) {
	filters := types.Map{
		"roomID": params.RoomID,
		"from": types.Map{
			"$gte": params.From,
		},
		"to": types.Map{
			"$lte": params.To,
		},
	}
	booking, err := h.store.Bookings.GetByFilter(ctx, filters)
	if err != nil {
		return false, err
	}
	if !booking.Cancelled {
		return true, err
	}
	return false, nil
}
