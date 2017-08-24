package model

import (
	"encoding/json"
	"strings"
)

// Visit of a user on certain place with rating
type Visit struct {
	ID         *int32    `json:"id,omitempty"`
	LocationID *int32    `json:"location"`
	Location   *Location `json:"-"`
	UserID     *int32    `json:"user"`
	User       *User     `json:"-"`
	VisitedAt  *int64    `json:"visited_at"`
	Mark       *uint8    `json:"mark"`
}

// VisitArray is a list of visits
type VisitArray struct {
	Visits []Visit `json:"visits"`
}

// UnmarshalJSON custom unmarshaller for Visit that converts timestamp to time struct
func (v *Visit) UnmarshalJSON(data []byte) error {
	// hot fix: do not allow null as a value
	if strings.Contains(string(data), ": null") {
		return ErrNullField
	}

	type Alias Visit
	var aux = (*Alias)(v)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
