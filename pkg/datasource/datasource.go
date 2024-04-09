package datasource

type AvailableResp struct {
	// 是否可用
	Available bool `json:"available"`
	// 消息
	Message string `json:"message"`
}

// Datasource 数据源接口
type Datasource interface {
	// Available 校验数据源是否可用
	Available(connect map[string]interface{}) AvailableResp

	// Invoke 数据源接口调用
	Invoke()
}
