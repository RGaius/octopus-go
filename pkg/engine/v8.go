package engine

import (
	"github.com/RGaius/octopus/pkg/util"
	"github.com/sirupsen/logrus"
	v8 "rogchap.com/v8go"
)

type V8 struct {
	iso        *v8.Isolate
	InvokeFunc *v8.FunctionTemplate
}

func NewV8() *V8 {
	iso := v8.NewIsolate()
	return &V8{
		iso:        iso,
		InvokeFunc: invokeFunc(iso),
	}
}

func (v *V8) Invoke(script string, params map[string]interface{}) (interface{}, error) {
	// 计算script md5值
	md5Val, _ := util.Md5Val(script)
	filename := md5Val + ".js"
	global := v8.NewObjectTemplate(v.iso) // a template that represents a JS Object
	global.Set("InvokeFunc", v.InvokeFunc)
	globalParam := v8.NewObjectTemplate(v.iso)
	// 遍历params，添加到globalParam
	for k, v := range params {
		globalParam.Set(k, v)
	}
	global.Set("ParamMap", globalParam)
	ctx := v8.NewContext(v.iso, global)
	val, _ := ctx.RunScript(script, filename)
	return converter(val)
}

// 接口调用函数
func invokeFunc(iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args := info.Args()
		// 如果参数为空，则抛出异常
		if len(args) == 0 {
			return nil
		}
		logrus.Infof("invokeFunc args: %v", args)
		return nil
	})
}
