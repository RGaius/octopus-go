package datatype

import "github.com/spf13/cast"

type Double struct {
}

func (d *Double) convert(value interface{}) (interface{}, error) {
	return cast.ToFloat64(d), nil
}

func init() {
	Register("DOUBLE", &Double{})
}
