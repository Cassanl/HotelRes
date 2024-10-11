package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"hoteRes/api"
	"hoteRes/db"
	"hoteRes/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi      = "mongodb://localhost:27017"
	dbtestName = "hotel-res-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func setupEnv(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{UserStore: db.NewMongoUserStore(client, dbtestName)}
}

func TestUserApi(t *testing.T) {
	tdb := setupEnv(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := api.NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "test@test.com",
		FirstName: "Pol",
		LastName:  "O'Brian",
		Password:  "Superstrongpwd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)
	// if err != nil {
	// 	t.Error(err)
	// }
	t.Log(resp)
}
