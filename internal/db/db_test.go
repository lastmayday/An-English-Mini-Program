package db

import (
	"testing"
)

func TestDbConfig(t *testing.T) {
	if DbConfig.DbUrl == "" {
		t.Error("get db config error")
	}
}

func TestCreateDbConn(t *testing.T) {
	_, err := CreateDbConn()
	if err != nil {
		t.Error("create db connection failed")
	}
}
