package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId     primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomId     primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NbPersons  uint               `bson:"nbPersons" json:"nbPersons"`
	From       time.Time          `bson:"from" json:"from"`
	To         time.Time          `bson:"to" json:"to"`
	CanceledAt time.Time          `bson:"canceledAt" json:"canceledAt"`
}

type BookingParams struct {
	NbPersons uint      `json:"nbPersons"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (params *BookingParams) Validate() map[string]string {
	errs := map[string]string{}
	if params.From.Before(time.Now()) {
		errs["from"] = "invalid From date : book before today"
	}
	if params.To.Before(params.From) {
		errs["to"] = "invalid To date : end date before start date"
	}
	return errs
}

func NewBookingFromParams(params BookingParams, userID, roomID primitive.ObjectID) *Booking {
	return &Booking{
		UserId:    userID,
		RoomId:    roomID,
		NbPersons: params.NbPersons,
		From:      params.From,
		To:        params.From,
	}
}
