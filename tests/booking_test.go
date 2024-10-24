package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hoteRes/api"
	"hoteRes/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBookingApi(t *testing.T) {
	var (
		tdb            = SetupEnv(t)
		conf           = fiber.Config{ErrorHandler: api.GlobalErrorHandler}
		app            = fiber.New(conf)
		bookingHandler = api.NewBookingHandler(tdb.store)
		authHandler    = api.NewAuthHandler(tdb.store.Users)
	)
	defer tdb.Teardown(t)

	// userParams := types.CreateUserParams{
	// 	Email:     "test@test.com",
	// 	FirstName: "Pol",
	// 	LastName:  "O'Brian",
	// 	Password:  "Superstrongpwd",
	// }
	// userDB, err := types.NewUserFromParams(userParams)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// user, err := tdb.store.Users.Insert(context.Background(), userDB)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	app.Post("/auth", authHandler.HandleAuthenticate)
	app.Post("/bookings", bookingHandler.HandlePostBooking)

	authParams := types.AuthParams{
		Email:    "test@test.com",
		Password: "turbopolo",
	}
	b, err := json.Marshal(authParams)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	_ = req

	oid, err := primitive.ObjectIDFromHex("67189d749b7d68faf24fd084")
	if err != nil {
		t.Fatal(err)
	}
	bookingParams := types.PostBookingParams{
		RoomID:    oid,
		NbPersons: 1,
		From:      time.Now().Add(time.Hour * 24),
		To:        time.Now().Add(time.Hour * 48),
	}
	b, err = json.Marshal(bookingParams)
	if err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest("POST", "/bookings", bytes.NewReader(b))
	_ = req
	// req.Header.Set("X-Api-Token", api.CreateTokenFromUser(user))

	app.Get("/ping", func(c *fiber.Ctx) error {
		t.Log("MONGO IS THE ISSUE")
		return c.Status(http.StatusOK).JSON(types.GenericResponse{
			Kind: types.OkResp,
			Msg:  "OK MATE",
		})
	})
	testReq := httptest.NewRequest("GET", "/ping", bytes.NewReader(make([]byte, 10)))
	resp, err := app.Test(testReq)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
