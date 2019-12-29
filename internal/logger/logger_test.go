package logger

import (
	"testing"
)

func TestInitLog(t *testing.T) {
	logFilePath := "/tmp/test.log"
	log := InitLog(logFilePath)
	if log == nil {
		t.Error("init log failed")
	}
}
