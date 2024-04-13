package datatype

import (
	"github.com/spf13/cast"
)

// 参数列表
var strTypeList = [...]string{"CHAR", "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT"}

type Str struct {
}

func (s *Str) convert(value interface{}) (interface{}, error) {
	return cast.ToString(value), nil
}

func init() {
	s := &Str{}
	for _, t := range strTypeList {
		Register(t, s)
	}
}
