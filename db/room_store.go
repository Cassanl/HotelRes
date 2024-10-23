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

	GetByFilter(context.Context, types.Filter) (*types.Room, error)
	ListByFilter(context.Context, types.Filter) ([]*types.Room, error)
	Insert(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(c *mongo.Client, hs HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     c,
		coll:       c.Database(DBNAME).Collection(roomColl),
		HotelStore: hs,
	}
}

func (s *MongoRoomStore) GetByFilter(ctx context.Context, filters types.Filter) (*types.Room, error) {
	var room types.Room
	if err := s.coll.FindOne(ctx, filters).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *MongoRoomStore) ListByFilter(ctx context.Context, filters types.Filter) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
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
