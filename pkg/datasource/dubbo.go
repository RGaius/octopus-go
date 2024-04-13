package datasource

type Dubbo string

func (d *Dubbo) Available(connect map[string]interface{}) AvailableResp {
	//TODO implement me
	panic("implement me")
}

func (d *Dubbo) Invoke(*InvokeParam) (interface{}, error) {
	//TODO implement me
	panic("implement me")
	return 0, nil
}
