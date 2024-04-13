package datasource

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestPostgres_Available(t *testing.T) {
	connect := make(map[string]interface{})
	connect["urlTemplate"] = "user={{.user}} password={{.password}} host={{.host}} port={{.port}} dbname={{.database}} sslmode=disable"
	connect["user"] = "synthetic_user"
	connect["password"] = "A0gqdr7c(Lfsuid"
	connect["host"] = "10.0.6.251"
	connect["port"] = "18104"
	connect["database"] = "cw_doau"
	t.Run("连接测试", func(t *testing.T) {
		p := &Postgres{}
		got := p.Available(connect)
		if !reflect.DeepEqual(got.Available, true) {
			t.Errorf("Available() = %v, want %v", got, true)
			return
		}
		t.Logf("msg:%s", got.Message)
	})
}

func TestPostgres_Invoke(t *testing.T) {
	connect := make(map[string]interface{})
	connect["urlTemplate"] = "user={{.user}} password={{.password}} host={{.host}} port={{.port}} dbname={{.database}} sslmode=disable"
	connect["user"] = "synthetic_user"
	connect["password"] = "A0gqdr7c(Lfsuid"
	connect["host"] = "10.0.6.251"
	connect["port"] = "18104"
	connect["database"] = "cw_doau"

	interfaceObj := make(map[string]interface{})
	interfaceObj["sql"] = "select * from aut_datasource"

	invokeParam := &InvokeParam{
		Name:       "test",
		Datasource: connect,
		Interface:  interfaceObj,
	}

	t.Run("执行测试", func(t *testing.T) {
		p := &Postgres{
			pool: make(map[string]*sql.DB),
		}
		_, _ = p.Invoke(invokeParam)
	})

}
