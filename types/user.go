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

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"-"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

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

func (params CreateUserParams) Validate() map[string]string {
	errs := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errs["firstName"] = fmt.Sprintf("invalid FirstName : required lenth is  %d", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errs["lastName"] = fmt.Sprintf("invalid LastName : required lenth is  %d", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errs["password"] = fmt.Sprintf("invalid Password : required lenth is  %d", minPasswordLen)
	}
	if _, err := mail.ParseAddress(params.Email); err != nil {
		errs["email"] = "invalid email address"
	}
	return errs
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params UpdateUserParams) ToFilter() Map {
	filters := Map{}
	if len(params.FirstName) > 0 {
		filters["firstName"] = params.FirstName
	}
	if len(params.LastName) > 0 {
		filters["lastName"] = params.LastName
	}
	return filters
}
