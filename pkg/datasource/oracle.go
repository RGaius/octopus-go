package datasource

import (
	"database/sql"
	"strings"
	"text/template"
)

type Oracle string

func (o *Oracle) Available(connect map[string]interface{}) AvailableResp {
	urlTemplate := connect["urlTemplate"]
	temple := template.Must(template.New("oracleUrl").Parse(urlTemplate.(string)))
	var result strings.Builder
	err := temple.Execute(&result, connect)
	if err != nil {
		return AvailableResp{
			Available: false,
			Message:   err.Error(),
		}
	}
	url := result.String()
	db, err := sql.Open("oracle", url)
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
	err = db.QueryRow("SELECT * FROM v$version").Scan(&version)
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

func (o *Oracle) Invoke(*InvokeParam) (interface{}, error) {
	//TODO implement me
	panic("implement me")
	return 0, nil
}
