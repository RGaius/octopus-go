package datasource

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/cast"
	"log"
	"strings"
	"sync"
	"text/template"
)

type Postgres struct {
	sync.Mutex // 添加互斥锁以保护数据库连接池的并发访问
	pool       map[string]*sql.DB
}

func (p *Postgres) Available(connect map[string]interface{}) AvailableResp {
	urlTemplate := connect["urlTemplate"]
	temple := template.Must(template.New("pgUrl").Parse(urlTemplate.(string)))
	var result strings.Builder
	err := temple.Execute(&result, connect)
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	url := result.String()

	db, err := sql.Open("postgres", url)
	defer db.Close()
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	err = db.Ping()
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
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

func (p *Postgres) Invoke(param *InvokeParam) (interface{}, error) {
	db, err := p.getOrCreate(param.Name, param.Datasource)
	if err != nil {
		log.Printf("Error getting or creating database connection: %v", err)
		return 0, nil
	}
	querySql := cast.ToString(param.Interface["sql"])
	log.Printf("执行postgre sql语句:%s", querySql)
	// 执行SQL语句，如果列表，则返回list map
	rows, err := db.Query(querySql)
	if err != nil {
		return nil, err
	}
	return ProcessRows(rows)
}

func (p *Postgres) getOrCreate(name string, datasource map[string]interface{}) (*sql.DB, error) {
	// 使用互斥锁保护pool的访问
	p.Lock()
	defer p.Unlock()
	db, ok := p.pool[name]
	if !ok {
		url, err := formatUrl(datasource)
		if err != nil {
			return nil, err // 优化错误处理，返回错误而不是直接返回nil
		}
		newDb, err := CreateDb("postgres", url)
		if err != nil {
			return nil, err // 处理createDb可能的错误
		}
		db = newDb
		p.pool[name] = db
	}
	return db, nil
}
