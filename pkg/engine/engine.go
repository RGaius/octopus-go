package engine

type Engine interface {
	Invoke(script string, params map[string]interface{}) (interface{}, error)
}
