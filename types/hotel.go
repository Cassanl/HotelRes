package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type UpdateHotelParams struct {
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

func (params UpdateHotelParams) ToBson() bson.M {
	result := bson.M{}
	if len(params.Name) > 0 {
		result["name"] = params.Name
	}
	if len(params.Location) > 0 {
		result["location"] = params.Location
	}
	if len(params.Rooms) > 0 {
		result["rooms"] = params.Rooms
	}
	return result
}

type RoomKind int

const (
	_ RoomKind = iota
	SingleRoomKind
	TwinRoomKind
	SeaSideRoomType
	DeluxeRoomKind
)

// type RoomType struct {
// 	ID        primitive.ObjectID
// 	BasePrice float64
// 	Kind      RoomKind
// }

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Kind      RoomKind           `bson:"kind" json:"kind"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
