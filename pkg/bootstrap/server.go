package bootstrap

import (
	"fmt"
	"github.com/RGaius/octopus/pkg/config"
	"github.com/RGaius/octopus/pkg/server"
	"github.com/RGaius/octopus/pkg/store"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

var (
	ConfigFilePath = ""
)

// Start 服务启动
func Start(configPath string) {
	ConfigFilePath = configPath
	// 解析配置文件
	cfg, err := config.Load(configPath)
	if err != nil {
		logrus.Error("load config fail")
		return
	}
	c, err := yaml.Marshal(cfg)
	if err != nil {
		logrus.Error("config yaml marshal fail")
		return
	}
	_, _ = fmt.Println(string(c))
	// 获取数据库配置
	store.SetStoreConfig(&cfg.Store)

	_, err = store.Load()
	if err != nil {
		logrus.Error("load store fail")
		return
	}

	httpServer := server.NewHTTPServer()
	err = httpServer.RegisterRouter()
	if err != nil {
		logrus.WithError(err).Errorln("Register router failed")
		os.Exit(1)
	}
	logrus.Info("Successfully registered router")
	err = httpServer.Run()
	if err != nil {
		logrus.WithError(err).Errorln("Start server failed")
		os.Exit(1)
	}
	logrus.Info(" start success")

}

func parseConfDir(path string) string {
	slashIndex := strings.LastIndex(path, "/")
	if slashIndex == -1 {
		return "./"
	}
	return path[0 : slashIndex+1]
}
