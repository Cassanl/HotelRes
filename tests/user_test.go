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
	var (
		tdb         = SetupEnv(t)
		conf        = fiber.Config{ErrorHandler: api.GlobalErrorHandler}
		app         = fiber.New(conf)
		userHandler = api.NewUserHandler(tdb.store.Users)
	)
	defer tdb.Teardown(t)

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
