package datasource

type Redis string

func (r *Redis) Available(map[string]interface{}) AvailableResp {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) Invoke(*InvokeParam) (interface{}, error) {
	//TODO implement me
	panic("implement me")
	return 0, nil
}
