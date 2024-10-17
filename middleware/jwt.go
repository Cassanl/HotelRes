package middleware

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
			return fmt.Errorf("unauthorized")
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		exp := int64(claims["exp"].(float64))
		expUtc := time.Unix(exp, 0).UTC()
		nowUtc := time.Now().UTC()
		if nowUtc.After(expUtc) {
			return fmt.Errorf("unauthorized")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		c.Context().SetUserValue("user", user)
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
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
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
