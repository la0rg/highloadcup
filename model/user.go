package model

import "time"

// User profile
type User struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

// UserArray is a list of users
type UserArray struct {
	Users []User `json:"users"`
}
