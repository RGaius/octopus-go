package datasource

import (
	"database/sql"
	"github.com/RGaius/octopus/pkg/datasource/datatype"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strings"
	"text/template"
)

type AvailableResp struct {
	// 是否可用
	Available bool `json:"available"`
	// 消息
	Message string `json:"message"`
}

type InvokeParam struct {
	// 命名空间
	Namespace string `json:"namespace"`
	// 数据源名称
	Name string `json:"name"`
	// 数据源连接信息
	Datasource map[string]interface{} `json:"datasource"`
	// 接口信息
	Interface map[string]interface{} `json:"interface"`
	// 请求参数
	ReqParam map[string]interface{} `json:"reqParam"`
}

// Datasource 数据源接口
type Datasource interface {
	// Available 校验数据源是否可用
	Available(connect map[string]interface{}) AvailableResp

	// Invoke 数据源接口调用
	Invoke(param *InvokeParam) (interface{}, error)
}

// formatUrl 使用给定的连接信息模板化URL。
//
// 参数:
// connect - 一个map类型的参数，包含至少一个键"urTemplate"，其值为一个字符串模板。
//
// 返回值:
// 返回一个字符串，表示根据模板和提供的连接信息生成的URL。
// 如果在解析模板或执行模板时发生错误，将返回一个空字符串和相应的错误。
func formatUrl(connect map[string]interface{}) (string, error) {
	// 从connect中提取URL模板字符串
	urlTemplate := cast.ToString(connect["urlTemplate"])
	// 解析URL模板
	temple := template.Must(template.New("url").Parse(urlTemplate))
	var result strings.Builder
	// 使用连接信息执行模板，结果存储在result中
	err := temple.Execute(&result, connect)
	if err != nil {
		return "", err // 如果执行过程中有错误，返回空字符串和错误信息
	}
	url := result.String()
	return url, nil
}

func (m *MySQL) GetOrCreate(name string, datasource map[string]interface{}) (*sql.DB, error) {
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

// CreateDb 创建数据库对象
func CreateDb(driverName string, url string) (*sql.DB, error) {
	// 创建数据库连接池
	db, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	// 设置连接池参数
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	return db, nil
}

// 假设rows来自一个数据库查询，columns是列的名称切片
func ProcessRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	// 获取列名
	columns, _ := rows.Columns()
	// 列类型
	columnTypes, _ := rows.ColumnTypes()

	// 使用defer确保rows在函数返回前被关闭
	defer rows.Close()

	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))

		for i := range columns {
			values[i] = &values[i] // 使用指针接收查询结果
		}

		if err := rows.Scan(values...); err != nil {
			// 改为记录错误，而不是直接退出
			logrus.Errorf("Error scanning row: %v", err)
			continue // 继续处理下一行，而不是终止整个函数
		}

		for i, col := range columns {
			columnType := columnTypes[i]
			// 打印
			if colVal := values[i]; colVal != nil {
				val := datatype.ToGoTypeValue(columnType.DatabaseTypeName(), colVal)
				logrus.Infof("Column: %s, Type: %s, Value: %v", col, columnType.DatabaseTypeName(), val)
				row[col] = val
			} else {
				row[col] = nil
			}
		}
		results = append(results, row)
	}

	// 确保没有其他错误，例如来自rows.Close()
	if err := rows.Err(); err != nil {
		logrus.Errorf("Error occurred during rows iteration: %v", err)
		return nil, err
	}
	return results, nil
}
