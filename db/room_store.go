package db

import (
	"context"
	"fmt"
	"hoteRes/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper

	InsertHotel(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRoomStore(c *mongo.Client, dbName string) *MongoRoomStore {
	return &MongoRoomStore{
		client: c,
		coll:   c.Database(dbName).Collection(roomColl),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// update hotel with room id
	return room, err
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("[INFO] dropping 'rooms' collection")
	return s.coll.Drop(ctx)
}
