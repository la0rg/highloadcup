package model

import (
	"encoding/json"
	"strings"
	"time"
)

// Visit of a user on certain place with rating
type Visit struct {
	ID         *int32     `json:"id,omitempty"`
	LocationID *int32     `json:"location"`
	UserID     *int32     `json:"user"`
	VisitedAt  *time.Time `json:"visited_at"`
	Mark       *uint8     `json:"mark"`
}

// VisitArray is a list of visits
type VisitArray struct {
	Visits []Visit `json:"visits"`
}

// MarshalJSON custom marshaller for Visit that converts time struct to timestamp
func (v *Visit) MarshalJSON() ([]byte, error) {
	type AliasVisit Visit
	return json.Marshal(struct {
		*AliasVisit
		VisitedAt int64 `json:"visited_at"`
	}{
		AliasVisit: (*AliasVisit)(v),
		VisitedAt:  v.VisitedAt.Unix(),
	})
}

// UnmarshalJSON custom unmarshaller for Visit that converts timestamp to time struct
func (v *Visit) UnmarshalJSON(data []byte) error {
	// hot fix: do not allow null as a value
	if strings.Contains(string(data), ": null") {
		return ErrNullField
	}
	type AliasVisit Visit
	aux := &struct {
		VisitedAt *int64 `json:"visited_at"`
		*AliasVisit
	}{
		AliasVisit: (*AliasVisit)(v),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.VisitedAt != nil {
		timestamp := time.Unix(*(aux.VisitedAt), 0)
		v.VisitedAt = &timestamp
	}
	return nil
}
