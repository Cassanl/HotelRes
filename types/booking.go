package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomId    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NbPersons int                `bson:"nbPersons" json:"nbPersons"`
	From      time.Time          `bson:"from" json:"from"`
	To        time.Time          `bson:"to" json:"to"`
}
