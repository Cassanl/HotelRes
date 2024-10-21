package main

import (
	"context"
	"errors"
	"flag"
	"hoteRes/api"
	"hoteRes/db"
	"hoteRes/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var appConf = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "no match"})
		}
		if errors.Is(err, primitive.ErrInvalidHex) {
			return c.JSON(map[string]string{"error": "invalid ID"})
		}
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Api server's listen address")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewBookingStore(client)
		store        = &db.Store{
			Users:    userStore,
			Hotels:   hotelStore,
			Rooms:    roomStore,
			Bookings: bookingStore,
		}
		app            = fiber.New(appConf)
		baseRouter     = app.Group("/api")
		v1Router       = baseRouter.Group("/v1", middleware.JWTAuthentication(userStore))
		adminRouter    = v1Router.Group("/admin", middleware.AdminAuth)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	v1Router.Get("/ping", handlePing)

	// ---------
	baseRouter.Post("/auth", authHandler.HandleAuthenticate)

	v1Router.Get("/", userHandler.HandleGetUsers)
	v1Router.Get("/:id", userHandler.HandleGetUser)
	v1Router.Post("/", userHandler.HandlePostUser)

	v1Router.Get("/", hotelHandler.HandleGetHotels)
	v1Router.Get("/:id", hotelHandler.HandleGetHotel)
	v1Router.Get("/:id/rooms", hotelHandler.HandleGetRooms)

	v1Router.Post("/:id/book", roomHandler.HandleBooking)
	v1Router.Delete("/:id", roomHandler.HandleCancelBooking)

	v1Router.Get("/:id", bookingHandler.HandleGetCurrentUserBookings)
	adminRouter.Get("/", bookingHandler.HandleGetBookings)
	adminRouter.Get("/:id", bookingHandler.HandleGetBooking)
	// ---------

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}
