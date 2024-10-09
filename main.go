package main

import (
	"context"
	"flag"
	"fmt"
	"hoteRes/api"
	"hoteRes/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-res"
	usercoll = "users"
)

func main() {
	mongo, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	coll := mongo.Database(dbname).Collection(usercoll)

	user := types.User{
		FirstName: "bil",
		LastName:  "bilbil",
	}
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	var bil types.User
	if err := coll.FindOne(ctx, bson.M{}).Decode(&bil); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", bil)

	listenAddr := flag.String("listenAddr", ":5000", "Api server's listen address")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler()
	apiv1.Get("/ping", handlePing)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}
