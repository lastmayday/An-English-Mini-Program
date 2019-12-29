package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lastmayday/BBC-English/web/app/article"
	"github.com/lastmayday/BBC-English/web/app/config"
)

func main() {
	router := gin.Default()

	router.GET("/articles", article.GetArticleSummaryList)
	router.GET("/article/:id", article.GetArticleDetail)

	port := config.AppConfig.AppPort
	router.Run(port)
}
