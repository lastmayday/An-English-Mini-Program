package article

import (
	"github.com/gin-gonic/gin"
	"github.com/lastmayday/BBC-English/internal/logger"
	"github.com/lastmayday/BBC-English/web/app/config"
	"github.com/lastmayday/BBC-English/web/app/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

var log *logrus.Logger

func init() {
	date := time.Now().Format("2006-01-02")
	logFilePath := config.AppConfig.AppLogFile + "." + date
	log = logger.InitLog(logFilePath)
}

func GetArticleSummaryList(c *gin.Context) {
	p := c.DefaultQuery("p", "1")
	pageSize := config.AppConfig.ArticlePageSize
	page, err := strconv.Atoi(p)
	if err != nil {
		log.WithFields(logrus.Fields{"page": p, "err": err}).Error("request page format error")
		c.JSON(http.StatusOK, gin.H{"success": false, "articles": nil, "errorMessage": "page format error"})
		return
	}

	articles, err := model.QueryArticleSummaryByPage(page, pageSize)
	if err != nil {
		log.WithFields(logrus.Fields{"page": p, "err": err}).Error("query article list error")
		c.JSON(http.StatusOK, gin.H{"success": false, "articles": nil, "errorMessage": "query db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "articles": articles})
}

func GetArticleDetail(c *gin.Context) {
	articleId := c.Param("id")
	article, err := model.QueryArticleDetail(articleId)

	if err != nil {
		log.WithFields(logrus.Fields{"articleId": articleId, "err": err}).Error("query article detail error")
		c.JSON(http.StatusOK, gin.H{"success": false, "article": nil, "errorMessage": "query db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "aritlce": article})
}
