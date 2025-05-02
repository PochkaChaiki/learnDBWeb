package config

import (
	"reflect"
	"testing"
)

func TestMustLoad(t *testing.T) {
	cfg := MustLoad()
	if reflect.ValueOf(cfg.Databases).Kind() != reflect.Map {
		t.Errorf("databases is not map")
	}
}
