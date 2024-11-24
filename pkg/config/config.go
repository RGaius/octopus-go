package config

import (
	"errors"
	"fmt"
	"github.com/RGaius/octopus/pkg/store"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	// 数据库连接配置
	Store store.Config
}

func Load(filePath string) (*Config, error) {
	if filePath == "" {
		err := errors.New("invalid config file path")
		fmt.Printf("[ERROR] %v\n", err)
		return nil, err
	}
	fmt.Printf("[INFO] load config from %v\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return nil, err
	}
	// 关闭文件
	defer func() {
		_ = file.Close()
	}()
	//	读取文件内容
	buf, err := os.ReadFile(filePath)
	if nil != err {
		return nil, fmt.Errorf("read file %s error", filePath)
	}

	conf := &Config{}
	if err := parseYamlContent(string(buf), conf); nil != err {
		return nil, err
	}
	return conf, nil
}

// parseYamlContent 解析yaml文件
func parseYamlContent(content string, conf *Config) error {
	if err := yaml.Unmarshal([]byte(replaceEnv(content)), conf); nil != err {
		return fmt.Errorf("parse yaml %s error:%w", content, err)
	}
	return nil
}

// replaceEnv replace holder by env list
func replaceEnv(configContent string) string {
	return os.ExpandEnv(configContent)
}
