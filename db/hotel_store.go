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

	GetByFilter(context.Context, types.Map) (*types.Hotel, error)
	GetById(context.Context, string) (*types.Hotel, error)
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	ListHotels(context.Context) ([]*types.Hotel, error)
	Update(context.Context, string, types.UpdateHotelParams) error
	AddRoom(context.Context, primitive.ObjectID, types.Map) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) GetByFilter(ctx context.Context, filter types.Map) (*types.Hotel, error) {
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, filter).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelStore) GetById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
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
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, id string, updateValues types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.D{{Key: "$push", Value: updateValues.ToDBMap()}}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	return nil

}

func (s *MongoHotelStore) AddRoom(ctx context.Context, hotelID primitive.ObjectID, updateValues types.Map) error {
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
