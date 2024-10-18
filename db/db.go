package db

import "context"

const (
	DBNAME      = "hotel-res"
	DBURI       = "mongodb://localhost:27017"
	TEST_DBNAME = "hotel-res-test"
	TEST_DBURI  = "mongodb://localhost:27018"
)

type Dropper interface {
	Drop(context.Context) error
}

type Store struct {
	Users    UserStore
	Hotels   HotelStore
	Rooms    RoomStore
	Bookings BookingStore
}
