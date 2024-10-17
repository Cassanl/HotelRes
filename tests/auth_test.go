package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"hoteRes/api"
	"hoteRes/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func insertTestUser(tdb *TestDB, t *testing.T) *types.User {
	userParams := types.CreateUserParams{
		Email:     "test@test.com",
		FirstName: "Pol",
		LastName:  "O'Brian",
		Password:  "Superstrongpwd",
	}

	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		t.Fatal(err)
	}
	tdb.store.Users.Insert(context.Background(), user)
	return user
}

func TestAuthenticate(t *testing.T) {
	tdb := SetupEnv(t)
	defer tdb.Teardown(t)

	app := fiber.New()
	authHandler := api.NewAuthHandler(tdb.store.Users)
	app.Post("/", authHandler.HandleAuthenticate)

	_ = insertTestUser(tdb, t)

	authParams := types.AuthParams{
		Email:    "test@ŧest.com",
		Password: "Superstrongpwd",
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200 OK but got : %d", resp.StatusCode)
	}
}
