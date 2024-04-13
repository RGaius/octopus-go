package datatype

import "github.com/spf13/cast"

var intTypeList = [...]string{"INT", "INT2"}

type Int struct {
}

func (i *Int) convert(value interface{}) (interface{}, error) {
	return cast.ToInt32(value), nil
}

func init() {
	i := &Int{}
	for _, s := range intTypeList {
		Register(s, i)
	}
}
