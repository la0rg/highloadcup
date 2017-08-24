package model

import "fmt"

type Avg struct {
	Value FloatPrec5 `json:"avg"`
}

type FloatPrec5 float64

func (n FloatPrec5) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.5f", float64(n))), nil
}
