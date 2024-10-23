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
	client *mongo.Client
	store  db.Store
	ctx    = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewBookingStore(client)
	)

	store = db.Store{
		Users:    userStore,
		Hotels:   hotelStore,
		Rooms:    roomStore,
		Bookings: bookingStore,
	}

	if err := store.Hotels.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := store.Rooms.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := store.Users.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := store.Bookings.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	res, err := store.Hotels.Insert(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			HotelID: res.ID,
			Size:    "small",
			Seaside: false,
			Price:   99.9,
		},
		{
			HotelID: res.ID,
			Size:    "medium",
			Seaside: true,
			Price:   199.9,
		},
		{
			HotelID: res.ID,
			Size:    "deluxe",
			Seaside: true,
			Price:   299.9,
		},
	}

	for _, room := range rooms {
		_, err := store.Rooms.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(isAdmin bool, fname, lname, email, pwsd string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  pwsd,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	if _, err = store.Users.Insert(ctx, user); err != nil {
		log.Fatal(err)
	}

}

func main() {
	seedHotel("Lacrustine", "Valencia")
	seedHotel("Al'Franco", "Madrid")
	seedUser(true, "polo", "POLO", "test@test.com", "turbopolo")
	seedUser(false, "james", "JAMES", "james@james.com", "turbojames")
}
