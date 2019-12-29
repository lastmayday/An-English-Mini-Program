package config

import (
	"testing"
)

func TestAppConfig(t *testing.T) {
	if AppConfig.ArticlePageSize == 0 {
		t.Error("init AppConfig failed")
	}
}
