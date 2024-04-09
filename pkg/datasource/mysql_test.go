package datasource

import (
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
			pool: new(map[string]Client),
		}
		if got := s.Available(connect); !reflect.DeepEqual(got.Available, true) {
			t.Errorf("Available() = %v, want %v", got, true)
		}
	})
}
