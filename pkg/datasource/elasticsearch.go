package datasource

type Elasticsearch struct {
}

func (e *Elasticsearch) Available(connect map[string]interface{}) AvailableResp {
	return AvailableResp{
		Available: true,
		Message:   "",
	}
}

func (e *Elasticsearch) Invoke(param *InvokeParam) (interface{}, error) {
	return nil, nil
}
