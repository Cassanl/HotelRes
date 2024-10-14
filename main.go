package main

import (
	"context"
	"errors"
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
		app   = fiber.New(appConf)
		apiv1 = app.Group("/api/v1")

		userStore  = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore  = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
	)

	apiv1.Get("/ping", handlePing)

	registerUserEndpoints(apiv1, userHandler)
	registerHotelEndpoints(apiv1, hotelHandler)

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
}
