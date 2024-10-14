package main

import (
	"context"
	"hoteRes/db"
	"hoteRes/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	ctx        = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

	if err := hotelStore.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := roomStore.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	res, err := hotelStore.Insert(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			HotelID:   res.ID,
			Kind:      types.SingleRoomKind,
			BasePrice: 99.9,
		},
		{
			HotelID:   res.ID,
			Kind:      types.SeaSideRoomType,
			BasePrice: 199.9,
		},
		{
			HotelID:   res.ID,
			Kind:      types.DeluxeRoomKind,
			BasePrice: 199.9,
		},
	}

	for _, room := range rooms {
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	seedHotel("Lacrustine", "Valencia")
	seedHotel("Al'Franco", "Madrid")
}
