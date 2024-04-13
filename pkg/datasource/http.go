package datasource

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type HTTP string

func (H *HTTP) Available(map[string]interface{}) AvailableResp {
	//TODO implement me
	panic("implement me")
}

func (H *HTTP) Invoke(param *InvokeParam) (interface{}, error) {
	datasource := param.Datasource
	// 获取host
	host := cast.ToString(datasource["host"])
	// 请求头
	datasourceHeaders := cast.ToStringMapString(datasource["headers"])

	// 从接口中获取请求方法
	interfaceObj := param.Interface
	// 获取url
	url := cast.ToString(interfaceObj["url"])
	// 获取请求头
	method := cast.ToString(interfaceObj["method"])
	// 获取接口请求头；与数据源请求头进行合并
	interfaceHeaders := cast.ToStringMapString(interfaceObj["headers"])
	// 合并datasourceHeaders与interfaceHeaders，若存在相同参数，则以interfaceHeaders为准
	err := mapstructure.Decode(interfaceHeaders, &datasourceHeaders)
	if err != nil {
		log.Fatalf("mapstructure.Decode err: %v", err)
	}
	// 使用fasthttp 构造请求
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(host + url)
	req.Header.SetMethod(method)
	for k, v := range interfaceHeaders {
		req.Header.Set(k, v)
	}
	// 设置一分钟超时时间
	req.SetTimeout(time.Minute)
	// 如果请求方法不为GET，则获取获取contentType
	if method != "GET" {
		contentType := cast.ToString(interfaceObj["contentType"])
		req.Header.SetContentType(contentType)
		// 如果contentType为application/json，则获取body，此时body为string类型
		if contentType == "application/json" {
			body := cast.ToString(interfaceObj["body"])
			req.SetBodyString(body)
		}
		// 如果contentType为application/x-www-form-urlencoded, 则获取body，此时body为数组，每个元素为键值对
		if contentType == "application/x-www-form-urlencoded" {
			body := cast.ToStringMapStringSlice(interfaceObj["body"])
			for k, v := range body {
				for _, vv := range v {
					req.PostArgs().Add(k, vv)
				}
			}
		}
		// 如果contentType为multipart/form-data, 则获取body，此时body为数组，每个元素为键值对
		if contentType == "multipart/form-data" {
			body := cast.ToStringMapStringSlice(interfaceObj["body"])
			for k, v := range body {
				for _, vv := range v {
					req.PostArgs().Add(k, vv)
				}
			}
		}
	}
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println(err)
		return nil, err
	}
	responseBody := cast.ToString(resp.Body())
	// 使用json格式化工具进行数据提取
	val, ok := interfaceObj["resultExtract"]
	if ok {
		return jsonpath.JsonPathLookup(responseBody, cast.ToString(val))
	}
	return jsonpath.JsonPathLookup(responseBody, "$")

}
