package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID      primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NbPersons   uint               `bson:"nbPersons" json:"nbPersons"`
	From        time.Time          `bson:"from" json:"from"`
	To          time.Time          `bson:"to" json:"to"`
	Cancelled   bool               `bson:"cancelled" json:"cancelled"`
	CancelledAt time.Time          `bson:"cancelledAt" json:"cancelledAt"`
}

type PostBookingParams struct {
	RoomID    primitive.ObjectID `json:"roomID"`
	NbPersons uint               `json:"nbPersons"`
	From      time.Time          `json:"from"`
	To        time.Time          `json:"to"`
}

func (params *PostBookingParams) Validate() map[string]string {
	errs := map[string]string{}
	if params.From.Before(time.Now()) {
		errs["from"] = "invalid From date : book before today"
	}
	if params.To.Before(params.From) {
		errs["to"] = "invalid To date : end date before start date"
	}
	return errs
}

func NewBookingFromParams(params PostBookingParams, userID primitive.ObjectID) *Booking {
	return &Booking{
		UserID:    userID,
		RoomID:    params.RoomID,
		NbPersons: params.NbPersons,
		From:      params.From,
		To:        params.From,
	}
}

// type UpdateBookingParams struct {
// 	NbPersons uint      `json:"nbPersons"`
// 	From      time.Time `json:"from"`
// 	To        time.Time `json:"to"`
// }

// func (params UpdateBookingParams) Validate() map[string]string {
// 	errs := map[string]string{}
// 	if params.From.Before(time.Now()) {
// 		errs["from"] = "invalid From date : book before today"
// 	}
// 	if params.To.Before(params.From) {
// 		errs["to"] = "invalid To date : end date before start date"
// 	}
// 	return errs
// }

// func (params UpdateBookingParams) ToFilter() Filter {

// }
