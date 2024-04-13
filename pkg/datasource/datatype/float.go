package datatype

import "github.com/spf13/cast"

type Float struct {
}

func (f *Float) convert(value interface{}) (interface{}, error) {
	return cast.ToFloat32(value), nil
}

func init() {
	Register("FLOAT", &Float{})
}
