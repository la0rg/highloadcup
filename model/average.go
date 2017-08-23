package model

import "strconv"

type Avg struct {
	Value FloatPrec5 `json:"avg"`
}

type FloatPrec5 float64

func (n FloatPrec5) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(n), 'f', -1, 64)), nil
}
