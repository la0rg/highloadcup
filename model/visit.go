package model

import "github.com/mailru/easyjson/opt"

// Visit of a user on certain place with rating
type Visit struct {
	ID         opt.Int32 `json:"id,omitempty"`
	LocationID opt.Int32 `json:"location"`
	Location   *Location `json:"-"`
	UserID     opt.Int32 `json:"user"`
	User       *User     `json:"-"`
	VisitedAt  opt.Int64 `json:"visited_at"`
	Mark       opt.Uint8 `json:"mark"`
}

// VisitArray is a list of visits
type VisitArray struct {
	Visits []Visit `json:"visits"`
}
