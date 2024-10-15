package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hoteRes/api"
	"hoteRes/db"
	"hoteRes/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.store.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func setupEnv(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}
	// TODO env test
	return &testdb{store: db.NewMongoUserStore(client)}
}

// TODO: currently broken
func TestUserApi(t *testing.T) {
	tdb := setupEnv(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := api.NewUserHandler(tdb.store)
	app.Post("/api/v1/users", userHandler.HandlePostUser)

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

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
	fmt.Println(resp.Status)
}
