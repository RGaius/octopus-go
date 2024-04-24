package engine

import (
	"encoding/json"
	"errors"
	v8 "rogchap.com/v8go"
)

func converter(val *v8.Value) (interface{}, error) {
	if val.IsString() || val.IsStringObject() {
		return val.String(), nil
	}
	if val.IsBoolean() {
		return val.Boolean(), nil
	}
	if val.IsInt32() {
		return val.Int32(), nil
	}
	if val.IsUint32() {
		return val.Uint32(), nil
	}
	if val.IsBigInt() {
		return val.BigInt(), nil
	}
	if val.IsNumber() || val.IsNumberObject() {
		return val.Number(), nil
	}
	if val.IsDate() {
		return val.String(), nil
	}
	if val.IsArray() {
		arrayVal, err := val.Object().MarshalJSON()
		if err != nil {
			return nil, err
		}
		var r = make([]interface{}, 0)
		err = json.Unmarshal(arrayVal, &r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if val.IsObject() {
		// 先转成json字符串字节数组再转成json对象
		jsonVal, err := val.Object().MarshalJSON()
		if err != nil {
			return nil, err
		}
		var r = make(map[string]interface{})
		err = json.Unmarshal(jsonVal, &r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if val.IsNullOrUndefined() {
		return nil, nil
	}
	if val.IsSymbol() {
		return nil, errors.New("symbol type not supported")
	}
	if val.IsFunction() {
		return nil, errors.New("function type not supported")
	}
	// 未知类型
	return nil, errors.New("unknown type")
}
