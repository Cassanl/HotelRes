package db

import (
	"context"
	"fmt"
	"hoteRes/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	Dropper

	Insert(context.Context, *types.Booking) (*types.Booking, error)
	ListByFilter(context.Context, types.Map) ([]*types.Booking, error)
	GetByFilter(context.Context, types.Map) (*types.Booking, error)
	Delete(context.Context, string) error
	Update(context.Context, types.Map, types.Map) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(c *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) GetByFilter(ctx context.Context, filters types.Map) (*types.Booking, error) {
	var booking types.Booking
	if err := s.coll.FindOne(ctx, filters).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) ListByFilter(ctx context.Context, filters types.Map) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoBookingStore) Update(ctx context.Context, id types.Map, updateValues types.Map) error {
	update := bson.D{{Key: "$set", Value: updateValues}}
	if _, err := s.coll.UpdateOne(ctx, id, update); err != nil {
		return err
	}
	return nil
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("[INFO] dropping 'bookings' collection")
	return s.coll.Drop(ctx)
}
