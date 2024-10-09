package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"hoteRes/api"
	"hoteRes/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var withConf = fiber.New(fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return nil
	},
})

func main() {
	mongo, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// ctx := context.Background()
	// user := types.User{
	// 	FirstName: "sabrina",
	// 	LastName:  "SABRINA",
	// }
	// coll := mongo.Database(db.DBNAME).Collection("users")
	// coll.InsertOne(ctx, user)

	listenAddr := flag.String("listenAddr", ":5000", "Api server's listen address")
	flag.Parse()

	// app := fiber.New()
	app := withConf
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/ping", handlePing)

	userStore := db.NewMongoUserStore(mongo)
	userHandler := api.NewUserHandler(userStore)
	// registerUserEndpoints(apiv1, *userHandler)

	apiv1User := apiv1.Group("/users")
	apiv1User.Get("/", userHandler.HandleGetUsers)
	apiv1User.Get("/:id", userHandler.HandleGetUser)
	apiv1User.Post("/", userHandler.HandleInsertUser)

	app.Listen(*listenAddr)
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"ping": "ping"})
}

func registerUserEndpoints(router fiber.Router, userHandler api.UserHandler) {
	userRoutes := router.Group("/users")
	userRoutes.Get("/", userHandler.HandleGetUsers)
	userRoutes.Get("/:id", userHandler.HandleGetUser)
	userRoutes.Post("/", userHandler.HandleInsertUser)
}
