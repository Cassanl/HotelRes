package api

import (
	"context"
	"hoteRes/db"
	"hoteRes/types"
	"net/http"

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

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Rooms.ListByFilter(c.Context(), types.Filter{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBooking(c *fiber.Ctx) error {
	var bookingParams types.BookingParams
	if err := c.BodyParser(&bookingParams); err != nil {
		return err
	}
	if errs := bookingParams.Validate(); len(errs) > 0 {
		return c.JSON(errs)
	}

	roomOid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue(types.UserKey).(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(types.GenericResponse{
			Kind: types.ErrorResp,
			Msg:  "unauthorized",
		})
	}

	booking := types.NewBookingFromParams(bookingParams, user.ID, roomOid)

	booked, err := h.isBooked(roomOid, bookingParams, c.Context())
	if err != nil {
		return err
	}
	if booked {
		return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{
			Kind: types.ErrorResp,
			Msg:  "already booked",
		})
	}
	res, err := h.store.Bookings.Insert(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// TODO
func (h *RoomHandler) HandleCancelBooking(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	var bookingParams types.BookingParams
	if err := c.BodyParser(&bookingParams); err != nil {
		return err
	}

	booked, err := h.isBooked(oid, bookingParams, c.Context())
	if err != nil {
		return err
	}
	if !booked {
		return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{
			Kind: types.ErrorResp,
			Msg:  "room is not booked : cannot cancel booking",
		})
	}

	if err := h.store.Bookings.Delete(c.Context(), oid.String()); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(types.GenericResponse{
		Kind: types.OkResp,
		Msg:  "can",
	})
}

func (h *RoomHandler) isBooked(oid primitive.ObjectID, params types.BookingParams, ctx context.Context) (bool, error) {
	filters := types.Filter{
		"roomID": oid,
		"from": types.Filter{
			"$gte": params.From,
		},
		"to": types.Filter{
			"$lte": params.To,
		},
	}
	bookings, err := h.store.Bookings.ListByFilter(ctx, filters)
	if err != nil {
		return false, err
	}
	if len(bookings) > 0 {
		return true, nil
	}
	return false, nil
}
