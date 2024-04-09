package datasource

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"strings"
)

type MySQL struct {
	// 系统连接池缓存
	pool *map[string]Client
}

type Client struct {
}

func (s *MySQL) Available(connect map[string]interface{}) AvailableResp {
	urlTemplate := connect["urlTemplate"]
	temple := template.Must(template.New("mysqlUrl").Parse(urlTemplate.(string)))
	var result strings.Builder
	err := temple.Execute(&result, connect)
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	url := result.String()

	db, err := sql.Open("mysql", url)
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

func (s *MySQL) Invoke() {

}
