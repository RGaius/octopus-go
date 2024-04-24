package engine

// Context 上下文
type Context struct {
	Source        string
	InterfaceList []string
	Params        map[string]interface{}
	Engine        Engine
}
