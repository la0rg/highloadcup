package model

import "github.com/mailru/easyjson/opt"

// Location of a sight
type Location struct {
	ID       opt.Int32  `json:"id"`
	Place    opt.String `json:"place"`
	Country  opt.String `json:"country"`
	City     opt.String `json:"city"`
	Distance opt.Int32  `json:"distance"`
}

// LocationArray is a list of locations
type LocationArray struct {
	Locations []Location `json:"locations"`
}

// // UnmarshalJSON custom unmarshaller for Location
// func (l *Location) UnmarshalJSON(data []byte) error {
// 	// hot fix: do not allow null as a value
// 	if strings.Contains(string(data), ": null") {
// 		return ErrNullField
// 	}
// 	type Alias Location
// 	var aux = (*Alias)(l)
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	return nil
// }
