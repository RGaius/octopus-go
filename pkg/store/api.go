package store

type Store interface {
	// Name 获取存储名称
	Name() string
	// Initialize 初始化
	Initialize(c *Config) error
	// AuthStore 认证
	AuthStore
}
