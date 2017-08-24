package model

import "encoding/json"

type UserVisit struct {
	Visit
}

type UserVisitArray struct {
	Visits []UserVisit `json:"visits"`
}

// MarshalJSON custom marshaller for UserVisit that converts time struct to timestamp
func (v *UserVisit) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Mark      uint8  `json:"mark"`
		Place     string `json:"place"`
		VisitedAt int64  `json:"visited_at"`
	}{
		Mark:      *(v.Mark),
		Place:     *(v.Location.Place),
		VisitedAt: *(v.VisitedAt),
	})
}
