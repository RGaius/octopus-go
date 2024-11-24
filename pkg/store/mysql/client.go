package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"time"
	"xorm.io/xorm"
)

type config struct {
	host             string
	port             string
	user             string
	password         string
	database         string
	maxOpenConns     int
	maxIdleConns     int
	connMaxLifetime  int
	txIsolationLevel int
}

// BaseDB 对sql.DB的封装
type BaseDB struct {
	*xorm.Engine
	cfg            *config
	isolationLevel sql.IsolationLevel
}

func NewBaseDB(cfg *config) (*BaseDB, error) {
	baseDb := &BaseDB{cfg: cfg}
	if cfg.txIsolationLevel > 0 {
		baseDb.isolationLevel = sql.IsolationLevel(cfg.txIsolationLevel)
		//log.Infof("[Store][database] use isolation level: %s", baseDb.isolationLevel.String())
	}
	if err := baseDb.openDatabase(); err != nil {
		return nil, err
	}
	return baseDb, nil
}

func (b *BaseDB) openDatabase() error {
	c := b.cfg

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.user, c.password, c.host, c.port, c.database)
	engine, err := xorm.NewEngine("mysql", dns)
	if err != nil {
		logrus.Errorf("[Store][database] sql open err: %s", err.Error())
		return err
	}
	if pingErr := engine.Ping(); pingErr != nil {
		logrus.Errorf("[Store][database] database ping err: %s", pingErr.Error())
		return pingErr
	}
	if c.maxOpenConns > 0 {
		logrus.Infof("[Store][database] db set max open conns: %d", c.maxOpenConns)
		engine.SetMaxOpenConns(c.maxOpenConns)
	}
	if c.maxIdleConns > 0 {
		logrus.Infof("[Store][database] db set max idle conns: %d", c.maxIdleConns)
		engine.SetMaxIdleConns(c.maxIdleConns)
	}
	if c.connMaxLifetime > 0 {
		logrus.Infof("[Store][database] db set conn max life time: %d", c.connMaxLifetime)
		engine.SetConnMaxLifetime(time.Second * time.Duration(c.connMaxLifetime))
	}
	b.Engine = engine
	return nil
}
