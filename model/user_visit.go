package model

type UserVisit struct {
	Mark      uint8  `json:"mark"`
	Place     string `json:"place"`
	VisitedAt int64  `json:"visited_at"`
}

func UserVisitFromVisit(v Visit) UserVisit {
	return UserVisit{
		Mark:      v.Mark.V,
		Place:     v.Location.Place.V,
		VisitedAt: v.VisitedAt.V,
	}
}

type UserVisitArray struct {
	Visits []UserVisit `json:"visits"`
}
