package store

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	Name    string
	Options map[string]interface{}
}

var (
	conf  = &Config{}
	Slots = make(map[string]Store)
	once  = &sync.Once{}
)

func SetStoreConfig(config *Config) {
	conf = config
}

// Load 加载存储层
func Load() (Store, error) {
	name := conf.Name
	store, ok := Slots[name]
	if !ok {
		logrus.Errorf("store `%s` not found", name)
		return nil, errors.New(fmt.Sprintf("store `%s` not found", name))
	}
	initialize(store)
	return store, nil
}

func initialize(s Store) {
	once.Do(func() {
		logrus.Infof("current use store plugin : %s\n", s.Name())
		if err := s.Initialize(conf); err != nil {
			logrus.Errorf("initialize store `%s` fail: %v", s.Name(), err)
			os.Exit(1)
		}
	})
}

// RegisterStore 注册一个新的Store
func RegisterStore(s Store) error {
	name := s.Name()
	if _, ok := Slots[name]; ok {
		return errors.New("store name already existed")
	}

	Slots[name] = s
	return nil
}
