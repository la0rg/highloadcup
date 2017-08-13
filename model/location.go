package model

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
