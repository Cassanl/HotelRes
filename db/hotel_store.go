package db

import (
	"context"
	"fmt"
	"hoteRes/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	Dropper

	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	ListHotels(context.Context) ([]*types.Hotel, error)
	Update(context.Context, string, types.UpdateHotelParams) error
	RefreshRooms(context.Context, primitive.ObjectID, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		client: c,
		coll:   c.Database(dbName).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, err
}

func (s *MongoHotelStore) ListHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		// TODO cannot decode array in objectid
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, id string, updateValues types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.D{{Key: "$push", Value: updateValues.ToBson()}}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	return nil

}

func (s *MongoHotelStore) RefreshRooms(ctx context.Context, hotelID primitive.ObjectID, updateValues bson.M) error {
	update := bson.D{{Key: "$push", Value: updateValues}}
	_, err := s.coll.UpdateByID(ctx, hotelID, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("[INFO] dropping 'hotels' collection")
	return s.coll.Drop(ctx)
}
