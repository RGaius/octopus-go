package datatype

import "github.com/spf13/cast"

type Mediumint struct {
}

func (m *Mediumint) convert(value interface{}) (interface{}, error) {
	return cast.ToInt32(value), nil
}

func init() {
	Register("MEDIUMINT", &Mediumint{})
}
