package types

import "golang.org/x/crypto/bcrypt"

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}
