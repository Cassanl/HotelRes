package api

import (
	"fmt"
	"hoteRes/db"
	"hoteRes/types"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if len(token) == 0 {
			return ErrUnauthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		exp := int64(claims["exp"].(float64))
		expUtc := time.Unix(exp, 0).UTC()
		nowUtc := time.Now().UTC()
		if nowUtc.After(expUtc) {
			return ErrUnauthorized()
		}

		userID := claims["id"].(string)
		user, err := userStore.GetById(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}

		c.Context().SetUserValue(types.UserKey, user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, ErrUnauthorized()
	}

	if !token.Valid {
		return nil, ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthorized()
	}

	return claims, nil
}

func CreateTokenFromUser(user *types.User) string {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 24).UTC().Unix(),
		"iat":   time.Now().UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	return tokenStr
}

func GetAuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue(types.UserKey).(*types.User)
	if !ok {
		return nil, ErrUnauthorized()
	}
	return user, nil
}
