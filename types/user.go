package types

import (
	"fmt"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errs := []string{}
	if len(params.FirstName) < minFirstNameLen {
		errs = append(errs, fmt.Sprintf("invalid FirstName : required lenth is  %d", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errs = append(errs, fmt.Sprintf("invalid LastName : required lenth is  %d", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errs = append(errs, fmt.Sprintf("invalid Password : required lenth is  %d", minPasswordLen))
	}
	if _, err := mail.ParseAddress(params.Email); err != nil {
		errs = append(errs, "invalid email address")
	}
	return errs
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

// type DbUser struct {
// 	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	FirstName         string             `bson:"firstName" json:"firstName"`
// 	LastName          string             `bson:"lastName" json:"lastName"`
// 	Email             string             `bson:"email" json:"email"`
// 	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
// 	CreatedAt         time.Time
// 	UpdatedAt         time.Time
// }

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
