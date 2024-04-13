package datatype

import "github.com/spf13/cast"

type Smallint struct {
}

func (s *Smallint) convert(value interface{}) (interface{}, error) {
	return cast.ToInt16(value), nil
}
func init() {
	Register("SMALLINT", &Smallint{})
}
