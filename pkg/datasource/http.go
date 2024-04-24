package datasource

import (
	"encoding/base64"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oliveagle/jsonpath"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
	"time"
)

type HTTP string

func (H *HTTP) Available(connection map[string]interface{}) AvailableResp {
	// 获取host
	host := cast.ToString(connection["host"])
	// 请求头
	datasourceHeaders := cast.ToStringMapString(connection["headers"])
	// 获取心跳地址
	heartbeat := cast.ToString(connection["heartbeat"])
	// 拼接地址
	url := fmt.Sprintf("%s%s", host, heartbeat)

	// 创建请求头
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	for k, v := range datasourceHeaders {
		req.Header.Set(k, v)
	}
	// 设置url
	req.SetRequestURI(url)
	// 设置请求方式
	req.Header.SetMethod("GET")
	// 发起http请求
	if err := fasthttp.DoTimeout(req, resp, 30*time.Second); err != nil {
		logrus.Errorf("http请求失败，err: %v", err)
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	// 若响应http状态不为200，则返回false
	if resp.StatusCode() != 200 {
		return AvailableResp{
			Available: false,
			Message:   cast.ToString(resp.Body()),
		}
	}
	return AvailableResp{
		Available: true,
		Message:   "success",
	}
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
		logrus.Errorf("mapstructure.Decode err: %v", err)
	}
	// 使用fasthttp 构造请求
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(host + url)
	req.Header.SetMethod(method)
	// 获取认证对象
	authObj := cast.ToStringMapString(datasource["authObj"])
	// 如果authObj不为空，则设置认证
	if authObj != nil && len(authObj) > 0 {
		wrapperAuth(datasource, authObj, req)
	}
	for k, v := range interfaceHeaders {
		req.Header.Set(k, v)
	}
	// 设置一分钟超时时间
	req.SetTimeout(time.Minute)
	// 如果请求方法不为GET，则获取获取contentType
	if method != "GET" {
		contentType := cast.ToString(interfaceObj["contentType"])
		req.Header.SetContentType(contentType)
		switch contentType {
		case "application/json":
			body := cast.ToString(interfaceObj["body"])
			req.SetBodyString(body)
		case "application/x-www-form-urlencoded":
			body := cast.ToStringMapStringSlice(interfaceObj["body"])
			for k, v := range body {
				for _, vv := range v {
					req.PostArgs().Add(k, vv)
				}
			}
		case "multipart/form-data":
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

// 封装认证
func wrapperAuth(datasource map[string]interface{}, authObj map[string]string, req *fasthttp.Request) {
	// 获取鉴权方式 basic
	authType := cast.ToString(datasource["authType"])
	if authType == "basic" {
		logrus.Infof("进行basic认证")
		// 获取用户名
		username := cast.ToString(authObj["username"])
		// 获取密码
		password := cast.ToString(authObj["password"])
		// 如果用户名或密码为空，则返回错误
		if username == "" || password == "" {
			logrus.Errorf("username or password is empty")
			return
		}
		// 生成basic认证
		req.Header.Set("Authorization", basicAuth(username, password))
		logrus.Infof("basic认证成功")
		return
	}
	// 普通认证，基于认证接口、请求方式、contentType及认证参数进行认证
	logrus.Infof("进行普通认证")
	url := cast.ToString(authObj["url"])
	method := cast.ToString(authObj["method"])
	if url == "" || method == "" {
		logrus.Errorf("url or method is empty")
		return
	}
	if method == "GET" {
		params, _ := cast.ToSliceE(authObj["params"])
		req.SetRequestURI(url)
		req.Header.SetMethod(method)
		//如果params不为nil或者params长度不为0，则遍历params
		if params != nil && len(params) != 0 {
			for _, v := range params {
				param := cast.ToStringMapString(v)
				// 获取参数名
				paramName := cast.ToString(param["name"])
				// 获取参数值
				paramValue := cast.ToString(param["value"])
				// 如果参数值为空，则返回错误
				if paramValue == "" {
					logrus.Errorf("param value is empty")
					return
				}
				// 设置url中的参数
				req.URI().QueryArgs().Add(paramName, paramValue)
			}
		}
		return
	}
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if method == "POST" {
		// 根据contentType进行组装数据
		contentType := cast.ToString(authObj["contentType"])
		switch contentType {
		case "application/json":
			// 获取body
			body := cast.ToString(authObj["body"])
			// 如果body为空，则返回错误
			if body == "" {
				logrus.Errorf("body is empty")
				return
			}
			req.SetBodyString(body)
		case "application/x-www-form-urlencoded":
			// 获取body
			body, _ := cast.ToSliceE(authObj["body"])
			for _, v := range body {
				param := cast.ToStringMapString(v)
				// 获取参数名
				paramName := cast.ToString(param["name"])
				// 获取参数值
				paramValue := cast.ToString(param["value"])
				// 如果参数值为空，则返回错误
				if paramValue == "" {
					logrus.Errorf("param value is empty")
					return
				}
				req.PostArgs().Add(paramName, paramValue)
			}
		case "multipart/form-data":
			// 获取body
			body, _ := cast.ToSliceE(authObj["body"])
			for _, v := range body {
				param := cast.ToStringMapString(v)
				// 获取参数名
				paramName := cast.ToString(param["name"])
				// 获取参数值
				paramValue := cast.ToString(param["value"])
				// 如果参数值为空，则返回错误
				if paramValue == "" {
					logrus.Errorf("param value is empty")
					return
				}
				req.PostArgs().Add(paramName, paramValue)
			}
		}
	}
}

// 基于用户名和密码生成basic认证
func basicAuth(username string, password string) string {
	authorization := username + ":" + password
	// 对authorization进行base64编码
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authorization))
}
