package main

import (
	"context"
	"errors"
	"flag"
	"hoteRes/api"
	"hoteRes/db"
	"hoteRes/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var withConf = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "no match"})
		}
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	mongo, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBNAME))
	if err != nil {
		log.Fatal(err)
	}

	// tempSeed(mongo)

	listenAddr := flag.String("listenAddr", ":5000", "Api server's listen address")
	flag.Parse()

	app := fiber.New(withConf)
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/ping", handlePing)

	userStore := db.NewMongoUserStore(mongo, db.DBNAME)
	userHandler := api.NewUserHandler(userStore)
	// registerUserEndpoints(apiv1, *userHandler)

	apiv1User := apiv1.Group("/users")
	apiv1User.Get("/", userHandler.HandleGetUsers)
	apiv1User.Get("/:id", userHandler.HandleGetUser)
	apiv1User.Post("/", userHandler.HandlePostUser)
	apiv1User.Put("/:id", userHandler.HandlePutUser)
	apiv1User.Delete("/:id", userHandler.HandleDeleteUser)

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}

// func registerUserEndpoints(router fiber.Router, userHandler api.UserHandler) {
// 	userRoutes := router.Group("/users")
// 	userRoutes.Get("/", userHandler.HandleGetUsers)
// 	userRoutes.Get("/:id", userHandler.HandleGetUser)
// 	userRoutes.Post("/", userHandler.HandlePostUser)
// }

func tempSeed(cl *mongo.Client) {
	ctx := context.Background()
	users := []types.User{
		{
			FirstName: "sabrina",
			LastName:  "SABRINA",
			Email:     "email@m",
		},
		{
			FirstName: "pol",
			LastName:  "POL",
			Email:     "email@pol",
		},
		{
			FirstName: "bil",
			LastName:  "BIL",
			Email:     "email@bil",
		},
		{
			FirstName: "heheh",
			LastName:  "jejejeA",
			Email:     "email@ddsdsds",
		},
		{
			FirstName: "non",
			LastName:  "NON",
			Email:     "email@nopn",
		},
	}
	coll := cl.Database(db.DBNAME).Collection("users")
	for _, user := range users {
		coll.InsertOne(ctx, user)
	}
}
