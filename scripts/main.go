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

const dburi = "mongodb://localhost:27017"

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)

	hotel := &types.Hotel{
		Name:     "SeasideAndalucia",
		Location: "Spain",
		Rooms:    []primitive.ObjectID{},
	}
	rooms := []types.Room{
		{
			Kind:      types.SingleRoomKind,
			BasePrice: 99.9,
		},
		{
			Kind:      types.SeaSideRoomType,
			BasePrice: 199.9,
		},
		{
			Kind:      types.DeluxeRoomKind,
			BasePrice: 199.9,
		},
	}

	resHotel, err := hotelStore.Insert(ctx, hotel)
	if err != nil {
		panic(err)
	}
	_ = resHotel

	for _, room := range rooms {
		res, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			panic(err)
		}
		_ = res
	}
}
