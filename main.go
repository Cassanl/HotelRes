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

var (
	client *mongo.Client
	app    *fiber.App
)

var withConf = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "no match"})
		}
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func init() {
	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	app = fiber.New(withConf)
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Api server's listen address")
	flag.Parse()

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/ping", handlePing)

	registerUserEndpoints(apiv1)
	registerHotelEndpoints(apiv1)

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}

func registerUserEndpoints(router fiber.Router) {
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	userHandler := api.NewUserHandler(userStore)

	userRoutes := router.Group("/users")

	userRoutes.Get("/", userHandler.HandleGetUsers)
	userRoutes.Get("/:id", userHandler.HandleGetUser)
	userRoutes.Post("/", userHandler.HandlePostUser)
}

func registerHotelEndpoints(router fiber.Router) {
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	hotelHandler := api.NewHotelHandler(hotelStore, roomStore)

	hotelRoutes := router.Group("/hotels")

	hotelRoutes.Get("/", hotelHandler.HandleGetHotels)
}
