package datatype

import "github.com/spf13/cast"

type Bigint struct {
}

func (b *Bigint) convert(value interface{}) (interface{}, error) {
	return cast.ToInt64(value), nil
}

func init() {
	Register("BIGINT", &Bigint{})
}
