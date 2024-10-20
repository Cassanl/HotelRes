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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var appConf = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "no match"})
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

	registerAuthEndpoint(baseRouter, authHandler)

	registerUserEndpoints(v1Router, userHandler)
	registerHotelEndpoints(v1Router, hotelHandler)
	registerRoomEndpoints(v1Router, roomHandler)

	registerBookingEndpoints(adminRouter, bookingHandler)

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}

func registerUserEndpoints(router fiber.Router, userHandler *api.UserHandler) {
	userRoutes := router.Group("/users")

	userRoutes.Get("/", userHandler.HandleGetUsers)
	userRoutes.Get("/:id", userHandler.HandleGetUser)
	userRoutes.Post("/", userHandler.HandlePostUser)
}

func registerHotelEndpoints(router fiber.Router, hotelHandler *api.HotelHandler) {
	hotelRoutes := router.Group("/hotels")

	hotelRoutes.Get("/", hotelHandler.HandleGetHotels)
	hotelRoutes.Get("/:id", hotelHandler.HandleGetHotel)
	hotelRoutes.Get("/:id/rooms", hotelHandler.HandleGetRooms)
}

func registerRoomEndpoints(router fiber.Router, roomHandler *api.RoomHandler) {
	roomRoutes := router.Group("/rooms")

	roomRoutes.Post("/:id/book", roomHandler.HandleBooking)
	roomRoutes.Delete("/:id/cancel", roomHandler.HandleCancelBooking)
}

func registerAuthEndpoint(router fiber.Router, authHandler *api.AuthHandler) {
	authRoutes := router.Group("/auth")

	authRoutes.Post("/", authHandler.HandleAuthenticate)
}

// TODO admin authz
func registerBookingEndpoints(router fiber.Router, bookingHandler *api.BookingHandler) {
	bookingRoutes := router.Group("/bookings")

	bookingRoutes.Get("/", bookingHandler.HandleGetBookings)
	bookingRoutes.Get("/:id", bookingHandler.HandleGetBooking)
}
