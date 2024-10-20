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
	Rating   int                  `bson:"rating" json:"rating"`
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

type HotelQueryParams struct {
	Rooms     bool
	MinRating int
	MaxRating int
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price   float64            `bson:"price" json:"price"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Size    string             `bson:"size" json:"size"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
