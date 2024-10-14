package db

import (
	"context"
	"fmt"
	"hoteRes/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(c *mongo.Client, dbName string, hs HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     c,
		coll:       c.Database(dbName).Collection(roomColl),
		HotelStore: hs,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	s.HotelStore.RefreshRooms(ctx, room.HotelID, bson.M{"rooms": room.ID})
	return room, err
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("[INFO] dropping 'rooms' collection")
	return s.coll.Drop(ctx)
}
