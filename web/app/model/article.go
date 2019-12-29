package model

import (
	"database/sql"
	database "github.com/lastmayday/BBC-English/internal/db"
)

type ArticleSummary struct {
	Id          string
	LastUpdated string
	Name        string
	Summary     string
}

type Article struct {
	Id          string
	LastUpdated string
	Name        string
	ShortName   string
	OriginUrl   string
	Body        string
	Type        string
	OssUrl      string
}

var db *sql.DB

func init() {
	db, _ = database.CreateDbConn()
}

func QueryArticleSummaryByPage(page int, pageSize int) ([]*ArticleSummary, error) {
	var articles []*ArticleSummary

	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT id, last_updated, name, summary FROM article ORDER BY gmt_create DESC LIMIT ?,?", offset, pageSize)
	if err != nil {
		return articles, err
	}
	defer rows.Close()

	for rows.Next() {
		var as ArticleSummary
		err = rows.Scan(&as.Id, &as.LastUpdated, &as.Name, &as.Summary)
		if err != nil {
			return articles, err
		}
		articles = append(articles, &as)
	}

	return articles, nil
}

func QueryArticleDetail(articleId string) (*Article, error) {
	var a Article

	stmt, err := db.Prepare("SELECT id, last_updated, name, short_name, origin_url, body, type, oss_url FROM article WHERE id = ?")
	if err != nil {
		return &Article{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(articleId).Scan(&a.Id, &a.LastUpdated, &a.Name, &a.ShortName, &a.OriginUrl, &a.Body, &a.Type, &a.OssUrl)
	return &a, err
}
