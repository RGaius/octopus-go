package engine

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestV8_Invoke(t *testing.T) {
	v := NewV8()
	params := make(map[string]interface{})
	params["a"] = 1

	str := "const arr=ParamMap['a'];arr"
	result, err := v.Invoke(str, params)
	if err != nil {
		logrus.Errorf("err:%v", err.Error())
	} else {
		logrus.Infof("result:%T, %v", result, result)
	}
}
