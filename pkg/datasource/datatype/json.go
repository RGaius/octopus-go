package datatype

import (
	"github.com/spf13/cast"
)

var jsonTypeList = [...]string{"JSON", "JSONB"}

type Json struct {
}

func (j *Json) convert(value interface{}) (interface{}, error) {
	return cast.ToString(value), nil
}

func init() {
	j := &Json{}
	for _, t := range jsonTypeList {
		Register(t, j)
	}
}
