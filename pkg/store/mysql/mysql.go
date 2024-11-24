package mysql

import (
	"fmt"
	"github.com/RGaius/octopus/pkg/store"
	"github.com/sirupsen/logrus"
)

// MySqlStore 实现了Store接口
type MySqlStore struct {
	start  bool
	baseDB *BaseDB
	*UserStore
}

const (
	// StoreName database storage name
	StoreName = "mysql"
	// DefaultConnMaxLifetime default maximum connection lifetime
	DefaultConnMaxLifetime = 60 * 30 // 默认是30分钟
	// emptyEnableTime 规则禁用时启用时间的默认值
	emptyEnableTime = "STR_TO_DATE('1980-01-01 00:00:01', '%Y-%m-%d %H:%i:%s')"
)

func init() {
	s := &MySqlStore{}
	_ = store.RegisterStore(s)
}

// Name 实现Name函数
func (s *MySqlStore) Name() string {
	return StoreName
}

func (s *MySqlStore) Initialize(c *store.Config) error {
	if s.start {
		return nil
	}
	// 解析数据库配置
	config, err := parseDatabaseConf(c.Options)
	if err != nil {
		return err
	}
	baseDB, err := NewBaseDB(config)
	if err != nil {
		return err
	}
	s.baseDB = baseDB
	logrus.Infof("[Store][database] connect the database successfully")

	s.start = true
	s.newStore()
	return nil
}

func (s *MySqlStore) newStore() {
	s.UserStore = &UserStore{
		Client: s.baseDB,
	}
}

func parseDatabaseConf(opt map[string]interface{}) (*config, error) {
	masterConfig, err := parseStoreConfig(opt)
	if err != nil {
		return nil, err
	}
	return masterConfig, nil
}

// parseStoreConfig 解析store的配置
func parseStoreConfig(opts map[string]interface{}) (*config, error) {
	needCheckFields := map[string]string{"user": "", "password": "", "host": "", "port": "", "database": ""}

	for key := range needCheckFields {
		val, ok := opts[key]
		if !ok {
			return nil, fmt.Errorf("config Plugin %s:%s type must be string", StoreName, key)
		}
		if val != nil {
			needCheckFields[key] = fmt.Sprintf("%v", val)
		} else {
			//log.Warnf("[Store][database] config field is empty: %s", key)
		}
	}

	c := &config{
		user:     needCheckFields["user"],
		password: needCheckFields["password"],
		host:     needCheckFields["host"],
		port:     needCheckFields["port"],
		database: needCheckFields["database"],
	}
	if maxOpenConns, _ := opts["maxOpenConns"].(int); maxOpenConns > 0 {
		c.maxOpenConns = maxOpenConns
	}
	if maxIdleConns, _ := opts["maxIdleConns"].(int); maxIdleConns > 0 {
		c.maxIdleConns = maxIdleConns
	}
	c.connMaxLifetime = DefaultConnMaxLifetime
	if connMaxLifetime, _ := opts["connMaxLifetime"].(int); connMaxLifetime > 0 {
		c.connMaxLifetime = connMaxLifetime
	}

	if isolationLevel, _ := opts["txIsolationLevel"].(int); isolationLevel > 0 {
		c.txIsolationLevel = isolationLevel
	}
	return c, nil
}
