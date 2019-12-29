package model

import (
	"testing"
)

func TestQueryArticleSummaryByPage(t *testing.T) {
	page := 1
	pageSize := 10
	articles, err := QueryArticleSummaryByPage(page, pageSize)
	if err != nil {
		t.Error("query article summary list failed")
	}
	if len(articles) == 0 {
		t.Error("query article summary length == 0")
	}
}

func TestQueryArticleDetail(t *testing.T) {
	articleId := "1"

	_, err := QueryArticleDetail(articleId)
	if err != nil {
		t.Error("query article detail failed")
	}
}
