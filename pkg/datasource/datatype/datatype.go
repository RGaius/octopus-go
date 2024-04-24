package datatype

import (
	"github.com/sirupsen/logrus"
)

type DataType interface {
	// convert value to the data type
	convert(value interface{}) (interface{}, error)
}

// 实现了DataType的类型集合
var interfaceMap = make(map[string]DataType)

// Register 注册函数，用于注册实现了DataType的类型
func Register(serviceName string, dataType DataType) {
	interfaceMap[serviceName] = dataType
}

// ToGoTypeValue 转换为Go类型
func ToGoTypeValue(columnType string, value interface{}) interface{} {
	// 如果 value为nil或者不存在类型，则返回value本身
	if value == nil {
		return value
	}
	if _, ok := interfaceMap[columnType]; !ok {
		logrus.Errorf("not found type:%s", columnType)
		return value
	}
	formatValue, err := interfaceMap[columnType].convert(value)
	if err != nil {
		return nil
	}
	return formatValue
}
