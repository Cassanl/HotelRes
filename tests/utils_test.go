package tests

import (
	"context"
	"hoteRes/db"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestDB struct {
	store *db.Store
}

func (tdb *TestDB) Teardown(t *testing.T) {
	if err := tdb.store.Users.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.store.Rooms.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.store.Hotels.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.store.Bookings.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func SetupEnv(t *testing.T) *TestDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewBookingStore(client)
	)

	return &TestDB{store: &db.Store{
		Users:    userStore,
		Hotels:   hotelStore,
		Rooms:    roomStore,
		Bookings: bookingStore,
	}}
}
