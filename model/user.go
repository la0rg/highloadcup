package model

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrNullField = errors.New("Field value is null")

// User profile
type User struct {
	ID        *int32  `json:"id"`
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Gender    *string `json:"gender"`
	BirthDate *int64  `json:"birth_date"`
}

// UserArray is a list of users
type UserArray struct {
	Users []User `json:"users"`
}

// UnmarshalJSON custom unmarshaller for User that converts timestamp to time struct
func (u *User) UnmarshalJSON(data []byte) error {
	// hot fix: do not allow null as a value
	if strings.Contains(string(data), ": null") {
		return ErrNullField
	}

	type Alias User
	var aux = (*Alias)(u)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
