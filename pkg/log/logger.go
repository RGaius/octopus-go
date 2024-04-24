package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

var (
	logger     *zap.Logger
	sugaredLog *zap.SugaredLogger
)

const defaultTimestampFormat = "2006-01-02T15:04:05.99999-07:00"

type LoggerConfig struct {
	Output     io.Writer
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"maxage"`
	MaxBackups int    `yaml:"maxbackups"`
	LocalTime  bool   `yaml:"localtime"`
	Compress   bool   `yaml:"compress"`
}

func InitLogger(config LoggerConfig) *logrus.Logger {
	logger := logrus.StandardLogger()
	if config.Output == nil {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		})
	} else {
		logger.SetOutput(config.Output)
	}

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		panic(fmt.Sprintf("parse log level: %+v", err))
	}
	logger.SetLevel(level)

	logger.SetReportCaller(true)
	return logger
}
