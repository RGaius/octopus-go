package datasource

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestMySQL_Available(t *testing.T) {
	connect := make(map[string]interface{})
	connect["urlTemplate"] = "{{.user}}:{{.password}}@tcp({{.host}}:{{.port}})/{{.database}}"
	connect["user"] = "Rootmaster"
	connect["password"] = "Rootmaster@777"
	connect["host"] = "10.2.2.153"
	connect["port"] = "18103"
	connect["database"] = "cw_doau"

	t.Run("连接测试", func(t *testing.T) {
		s := &MySQL{
			pool: make(map[string]*sql.DB),
		}
		if got := s.Available(connect); !reflect.DeepEqual(got.Available, true) {
			t.Errorf("Available() = %v, want %v", got, true)
		}
	})
}

func TestMySQL_Invoke(t *testing.T) {
	connect := make(map[string]interface{})
	connect["urlTemplate"] = "{{.user}}:{{.password}}@tcp({{.host}}:{{.port}})/{{.database}}"
	connect["user"] = "Rootmaster"
	connect["password"] = "Rootmaster@777"
	connect["host"] = "10.2.2.153"
	connect["port"] = "18103"
	connect["database"] = "cw_doau"

	interfaceObj := make(map[string]interface{})
	interfaceObj["sql"] = "select * from aut_datasource"

	invokeParam := &InvokeParam{
		Name:       "test",
		Datasource: connect,
		Interface:  interfaceObj,
	}

	t.Run("执行测试", func(t *testing.T) {
		s := &MySQL{
			pool: make(map[string]*sql.DB),
		}
		result, _ := s.Invoke(invokeParam)
		t.Logf("result:%v", result)
	})

}
