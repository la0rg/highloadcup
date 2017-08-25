package model

import (
	"errors"

	"github.com/mailru/easyjson/opt"
)

var ErrNullField = errors.New("Field value is null")

// User profile
type User struct {
	ID        opt.Int32  `json:"id"`
	Email     opt.String `json:"email"`
	FirstName opt.String `json:"first_name"`
	LastName  opt.String `json:"last_name"`
	Gender    opt.String `json:"gender"`
	BirthDate opt.Int64  `json:"birth_date"`
}

// UserArray is a list of users
type UserArray struct {
	Users []User `json:"users"`
}