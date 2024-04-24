package engine

import (
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type OttoEngine struct {
	otto otto.Otto
}

func NewOtto() *OttoEngine {
	o := otto.New()
	return &OttoEngine{otto: *o}
}

func (o *OttoEngine) Invoke(script string, params map[string]interface{}) (interface{}, error) {
	runtime := o.otto.Copy()
	runtime.Set("Param", params)
	runtime.Set("Invoke", func(call otto.FunctionCall) otto.Value {
		//argument := call.ArgumentList
		// 调用
		return otto.Value{}
	})
	value, err := runtime.Run(script)
	if err != nil {
		logrus.Errorf("OttoEngine invoke error: %v", err.Error())
		return nil, err
	}
	return value.Export()
}
