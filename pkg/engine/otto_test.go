package engine

import (
	"reflect"
	"testing"
)

func TestNewOtto(t *testing.T) {
	tests := []struct {
		name string
		want *OttoEngine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOtto(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOtto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOttoEngine_Invoke(t *testing.T) {
	o := NewOtto()
	params := make(map[string]interface{})
	params["1"] = 123
	script := `
	var aa = [1,2,3];aa
`

	got, err := o.Invoke(script, params)
	if err != nil {
		t.Errorf("Invoke() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, "123") {
		t.Errorf("Invoke() got = %v, want %v", got, "123")
	}
}
