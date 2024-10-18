package api

import (
	"fmt"
	"hoteRes/db"
	"hoteRes/types"
	"net/http"
	"slices"

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
			Kind:   types.ErrorResp,
			Status: fiber.StatusInternalServerError,
		})
	}

	booking := types.NewBookingFromParams(bookingParams, user.ID, roomOid)
	filters := types.Filter{
		"fromDate": types.Filter{
			"$gte": bookingParams.From,
		},
		"toDate": types.Filter{
			"$lte": bookingParams.To,
		},
	}
	bookings, err := h.store.Bookings.ListByFilter(c.Context(), filters)
	if err != nil {
		return err
	}
	fmt.Println(bookings)
	if ok := slices.Contains[[]*types.Booking](bookings, booking); ok {
		return c.JSON("room already booked")
	}
	res, err := h.store.Bookings.Insert(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(res)
}
