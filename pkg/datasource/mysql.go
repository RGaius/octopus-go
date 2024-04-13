package datasource

import (
	"database/sql"
	"errors"
	"github.com/RGaius/octopus/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
	"sync"
	"time"
)

const (
	maxOpenConns    = 10
	maxIdleConns    = 5
	connMaxLifetime = time.Minute * 3
)

type MySQL struct {
	sync.Mutex // 添加互斥锁以保护数据库连接池的并发访问
	pool       map[string]*sql.DB
}

func (m *MySQL) Available(connect map[string]interface{}) AvailableResp {
	url, err := formatUrl(connect)
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	db, err := sql.Open("mysql", url)
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	var version string
	available, err := executeQuery(db, "SELECT VERSION()", func(row *sql.Row) error {
		return row.Scan(&version)
	})
	if !available || err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	return AvailableResp{
		Available: true,
		Message:   version,
	}
}

func (m *MySQL) Invoke(param *InvokeParam) (interface{}, error) {
	db, err := m.getOrCreate(param.Name, param.Datasource)
	if err != nil {
		log.GetSugaredLogger().Warnf("Error getting or creating database connection:%v", err)
		return 0, nil
	}
	querySql := cast.ToString(param.Interface["sql"])
	log.GetSugaredLogger().Infof("sql:%s", querySql)
	// 执行SQL语句，如果列表，则返回list map
	rows, err := db.Query(querySql)
	if err != nil {
		return nil, err
	}
	return ProcessRows(rows)
}

// 将直接操作数据库的逻辑封装，减少重复代码，优化错误处理。
func executeQuery(db *sql.DB, query string, scanFunc func(*sql.Row) error) (bool, error) {
	row := db.QueryRow(query)
	err := scanFunc(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // 如果没有行，返回false，而不是错误
		}
		return false, err // 其他错误，返回false和错误信息
	}
	return true, nil // 扫描成功，返回true
}

func (m *MySQL) getOrCreate(name string, datasource map[string]interface{}) (*sql.DB, error) {
	// 使用互斥锁保护pool的访问
	m.Lock()
	defer m.Unlock()
	db, ok := m.pool[name]
	if !ok {
		url, err := formatUrl(datasource)
		if err != nil {
			return nil, err // 优化错误处理，返回错误而不是直接返回nil
		}
		newDb, err := CreateDb("mysql", url)
		if err != nil {
			return nil, err // 处理createDb可能的错误
		}
		db = newDb
		m.pool[name] = db
	}
	return db, nil
}
