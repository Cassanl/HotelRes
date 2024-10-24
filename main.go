package main

import (
	"context"
	"flag"
	"hoteRes/api"
	"hoteRes/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var appConf = fiber.Config{
	ErrorHandler: api.GlobalErrorHandler,
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
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Users:    userStore,
			Hotels:   hotelStore,
			Rooms:    roomStore,
			Bookings: bookingStore,
		}
		app            = fiber.New(appConf)
		baseRouter     = app.Group("/api")
		v1Router       = baseRouter.Group("/v1", api.JWTAuthentication(userStore))
		adminRouter    = v1Router.Group("/admin", api.AdminAuth)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	v1Router.Get("/ping", handlePing)

	// ---------
	baseRouter.Post("/auth", authHandler.HandleAuthenticate)

	v1Router.Get("/hotels", hotelHandler.HandleGetHotels)
	v1Router.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	// TODO generic filter
	// v1Router.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	v1Router.Get("/bookings", bookingHandler.HandleGetCurrentUserBookings)
	v1Router.Post("/bookings", bookingHandler.HandlePostBooking)
	v1Router.Delete("bookings/:id", bookingHandler.HandleCancelBooking)

	v1Router.Get("/rooms", roomHandler.HandleGetRooms)
	v1Router.Get("/rooms/:id", roomHandler.HandleGetRoom)

	adminRouter.Get("/users", userHandler.HandleGetUsers)
	adminRouter.Get("/users/:id", userHandler.HandleGetUser)
	adminRouter.Post("/users", userHandler.HandlePostUser)

	adminRouter.Get("/bookings", bookingHandler.HandleGetBookings)
	adminRouter.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	adminRouter.Delete("bookings/:id", bookingHandler.HandleDeleteBooking)
	// ---------

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}
