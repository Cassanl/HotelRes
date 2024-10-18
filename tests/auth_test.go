package tests

import (
	"hoteRes/api"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticate(t *testing.T) {
	tdb := SetupEnv(t)
	defer tdb.Teardown(t)

	app := fiber.New()
	authHandler := api.NewAuthHandler(tdb.store.Users)
	app.Post("/", authHandler.HandleAuthenticate)
}
