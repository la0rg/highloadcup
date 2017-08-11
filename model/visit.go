package model

import "time"

// Visit of a user on certain place with rating
type Visit struct {
	ID         int32     `json:"id,omitempty"`
	LocationID int32     `json:"location"`
	UserID     int32     `json:"user"`
	VisitedAt  time.Time `json:"visited_at"`
	Mark       uint8     `json:"mark"`
}

// VisitArray is a list of visits
type VisitArray struct {
	Visits []Visit `json:"visits"`
}
