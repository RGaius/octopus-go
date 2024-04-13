package datatype

import (
	"github.com/spf13/cast"
	"log"
	"time"
)

var timeTypeList = [...]string{"DATE", "TIMESTAMP", "TIME", "DATETIME", "YEAR"}

type Time struct {
}

func (t *Time) convert(v interface{}) (interface{}, error) {
	// 现将值转换为字符串
	strVal := cast.ToString(v)
	temp, err := cast.ToTimeE(strVal)
	if err != nil {
		log.Println("Error converting DATETIME:", err)
		return time.Now(), err
	}
	return temp, nil
}

func init() {
	t2 := &Time{}
	for _, t := range timeTypeList {
		Register(t, t2)
	}
}
