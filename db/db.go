package db

import "context"

const (
	DBNAME      = "hotel-res"
	TEST_DBNAME = "hotel-res-test"
	DBURI       = "mongodb://localhost:27017"
)

type Dropper interface {
	Drop(context.Context) error
}
