package tests

import (
	"bytes"
	"encoding/json"
	"hoteRes/api"
	"hoteRes/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUserApi(t *testing.T) {
	tdb := SetupEnv(t)
	defer tdb.Teardown(t)

	app := fiber.New()
	userHandler := api.NewUserHandler(tdb.store.Users)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "test@test.com",
		FirstName: "Pol",
		LastName:  "O'Brian",
		Password:  "Superstrongpwd",
	}

	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	// resp, err := app.Test(req)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(resp)
}
