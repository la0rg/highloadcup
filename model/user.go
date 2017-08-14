package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrNullField = errors.New("Field value is null")

// User profile
type User struct {
	ID        *int32     `json:"id"`
	Email     *string    `json:"email"`
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Gender    *string    `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
}

// UserArray is a list of users
type UserArray struct {
	Users []User `json:"users"`
}

// MarshalJSON custom marshaller for User that converts time struct to timestamp
func (u *User) MarshalJSON() ([]byte, error) {
	type AliasUser User
	return json.Marshal(struct {
		*AliasUser
		BirthDate int64 `json:"birth_date"`
	}{
		AliasUser: (*AliasUser)(u),
		BirthDate: u.BirthDate.Unix(),
	})
}

// UnmarshalJSON custom unmarshaller for User that converts timestamp to time struct
func (u *User) UnmarshalJSON(data []byte) error {
	// hot fix: do not allow null as a value
	if strings.Contains(string(data), ": null") {
		return ErrNullField
	}

	type AliasUser User
	aux := &struct {
		BirthDate *int64 `json:"birth_date"`
		*AliasUser
	}{
		AliasUser: (*AliasUser)(u),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.BirthDate != nil {
		timestamp := time.Unix(*(aux.BirthDate), 0)
		u.BirthDate = &timestamp
	}
	return nil
}
