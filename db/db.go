package db

import "context"

const DBNAME = "hotel-res"

type Dropper interface {
	Drop(context.Context) error
}
