package log

import "go.uber.org/zap"

var (
	logger     *zap.Logger
	sugaredLog *zap.SugaredLogger
)

func init() {
	// 根据实际情况配置日志级别、编码器、输出等
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.InfoLevel) // 示例：设置日志级别为 Info
	config.Encoding = "console"          // 示例：使用颜色丰富的控制台输出

	// 创建并初始化 Logger
	var err error
	logger, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// 创建 SugaredLogger
	sugaredLog = logger.Sugar()
}

// GetLogger 返回全局的 zap.Logger 实例
func GetLogger() *zap.Logger {
	return logger
}

// GetSugaredLogger 返回全局的 zap.SugaredLogger 实例
func GetSugaredLogger() *zap.SugaredLogger {
	return sugaredLog
}
