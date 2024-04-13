package datatype

import "github.com/spf13/cast"

var tinyintTypeList = [...]string{"TINYINT", "INT8"}

type Tinyint struct {
}

func (t *Tinyint) convert(value interface{}) (interface{}, error) {
	return cast.ToInt8(value), nil
}

func init() {
	t := &Tinyint{}
	for _, s := range tinyintTypeList {
		Register(s, t)
	}
}
