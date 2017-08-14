package model

import (
	"encoding/json"
	"strings"
)

// Location of a sight
type Location struct {
	ID       *int32  `json:"id"`
	Place    *string `json:"place"`
	Country  *string `json:"country"`
	City     *string `json:"city"`
	Distance *int32  `json:"distance"`
}

// LocationArray is a list of locations
type LocationArray struct {
	Locations []Location `json:"locations"`
}

// UnmarshalJSON custom unmarshaller for Location
func (l *Location) UnmarshalJSON(data []byte) error {
	// hot fix: do not allow null as a value
	if strings.Contains(string(data), ": null") {
		return ErrNullField
	}
	type Alias Location
	var aux = (*Alias)(l)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
