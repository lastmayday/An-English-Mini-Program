package bbc

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	database "github.com/lastmayday/BBC-English/internal/db"
	"github.com/lastmayday/BBC-English/internal/logger"
	"github.com/lastmayday/BBC-English/internal/oss"
	"github.com/lastmayday/BBC-English/internal/xml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// -------------- BBC API struct start ---------------

type BBCHome struct {
	Relations []HomeRelation
}

type HomeRelation struct {
	Content HomeContent
}

type HomeContent struct {
	Id string
}

type BBCPost struct {
	Id          string
	LastUpdated int64
	Name        string
	Summary     string
	ShortName   string
	Body        string
	Relations   []MediaRelation
	ShareUrl    string
}

type MediaRelation struct {
	PrimaryType string
	Content     MediaContent
}

type MediaContent struct {
	Id         string
	Type       string
	ExternalId string
}

type BBCMedia struct {
	Media []ExternalMedia
}

type ExternalMedia struct {
	Bitrate    string
	Type       string
	Kind       string
	Service    string
	Connection []MediaConnection
}

type MediaConnection struct {
	Protocol       string
	Href           string
	AuthExpires    string
	TransferFormat string
	Priority       string
	Supplier       string
}

// -------------- BBC API struct end ---------------

// -------------- database struct start ---------------

type Article struct {
	BBCId       string
	LastUpdated string
	Name        string
	Summary     string
	ShortName   string
	OriginUrl   string
	Body        string
	MediaId     string
	Type        string
	OssUrl      string
	FilePath    string
	FileName    string
}

// -------------- database struct end ---------------

// -------------- config struct start ---------------

type BBCConfig struct {
	BaseContentUrl  string `mapstructure:"BASE_CONTENT_URL"`
	HomePath        string `mapstructure:"HOME_PATH"`
	BaseMediaUrl    string `mapstructure:"BASE_MEDIA_URL"`
	SuffixMediaPath string `mapstructure:"SUFFIX_MEDIA_PATH"`
	DownloadPath    string `mapstructure:"DOWNLOAD_PATH"`
	CrawlerLogFile  string `mapstructure:"CRAWLER_LOG_FILE"`
}

// -------------- config struct end ---------------

var bbcConfig BBCConfig
var httpClient = &http.Client{Timeout: 10 * time.Second}
var uploadChannel = make(chan *Article)
var log *logrus.Logger

func init() {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("bbc_config")
	v.AddConfigPath("$HOME/go/src/github.com/lastmayday/BBC-English/config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&bbcConfig)
	if err != nil {
		panic(err)
	}

	date := time.Now().Format("2006-01-02")
	logFilePath := bbcConfig.CrawlerLogFile + "." + date
	log = logger.InitLog(logFilePath)
}

func crawl(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func crawlHome() *BBCHome {
	log.Info("start to crawl home...")
	home := new(BBCHome)
	crawl(bbcConfig.BaseContentUrl+bbcConfig.HomePath, home)
	log.Info("crawl home done!")
	return home
}

func crawlPost(postId string) *BBCPost {
	log.WithFields(logrus.Fields{"postId": postId}).Info("start to crawl post")
	postUrl := bbcConfig.BaseContentUrl + postId
	bbcPost := new(BBCPost)
	crawl(postUrl, bbcPost)
	log.WithFields(logrus.Fields{"postId": postId}).Info("crawl post done!")
	return bbcPost
}

func crawlMedia(mediaId string) *BBCMedia {
	log.WithFields(logrus.Fields{"mediaId": mediaId}).Info("start to crawl media")
	mediaUrl := bbcConfig.BaseMediaUrl + mediaId + bbcConfig.SuffixMediaPath
	bbcMedia := new(BBCMedia)
	crawl(mediaUrl, bbcMedia)
	log.WithFields(logrus.Fields{"mediaId": mediaId}).Info("crawl media done!")
	return bbcMedia
}

func generateArticle(post *BBCPost) *Article {
	article := new(Article)
	article.BBCId = post.Id
	article.LastUpdated = strconv.FormatInt(post.LastUpdated, 10)
	article.Name = post.Name
	article.Summary = post.Summary
	article.ShortName = post.ShortName
	article.OriginUrl = post.ShareUrl
	article.Body = xml.XmlToHtml(post.Body)
	return article
}

func executeCommand(cmdName string, cmdArgs []string) {
	log.WithFields(logrus.Fields{"cmdName": cmdName, "cmdArgs": cmdArgs}).Info("cmd execution start")
	_, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		log.WithFields(logrus.Fields{"cmdName": cmdName, "cmdArgs": cmdArgs, "err": err}).Error("cmd execute error")
	}
	log.WithFields(logrus.Fields{"cmdName": cmdName, "cmdArgs": cmdArgs}).Info("cmd execution end")
}

func downloadAndUpload(article *Article, mediaUrl string) {
	cmdName := "youtube-dl"
	cmdArgs := []string{"-o", article.FilePath, mediaUrl}

	executeCommand(cmdName, cmdArgs)

	ossUrl, err := oss.Upload(article.FileName, article.FilePath)
	if err != nil {
		log.WithFields(logrus.Fields{"filePath": article.FilePath, "err": err}).Error("upload failed")
		panic(err)
	} else {
		log.WithFields(logrus.Fields{"fileName": article.FileName}).Info("uploaded")
		article.OssUrl = ossUrl

		rmCmdName := "rm"
		rmCmdArgs := []string{"-fr", article.FilePath}
		executeCommand(rmCmdName, rmCmdArgs)
	}
	uploadChannel <- article
}

func isArticleExsited(db *sql.DB, bbcId string) bool {
	stmt, err := db.Prepare("SELECT id FROM article WHERE bbc_id = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var id sql.NullString
	err = stmt.QueryRow(bbcId).Scan(&id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false
		default:
			panic(err)
		}
	}
	return id.Valid
}

func insertDb(db *sql.DB, article *Article) {
	stmt, err := db.Prepare(`
	INSERT INTO article(bbc_id, last_updated, name, summary, short_name,
	origin_url, body, media_id, type, oss_url)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.BBCId, article.LastUpdated, article.Name, article.Summary, article.ShortName,
		article.OriginUrl, article.Body, article.MediaId, article.Type, article.OssUrl)
	if err != nil {
		panic(err)
	}

	log.WithFields(logrus.Fields{"bbcId": article.BBCId}).Info("insert success")
}

func CrawlArticles() {
	db, err := database.CreateDbConn()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	home := crawlHome()

	executeCount := 0
	for _, relation := range home.Relations {
		exsited := isArticleExsited(db, relation.Content.Id)
		if exsited {
			log.WithFields(logrus.Fields{"bbcId": relation.Content.Id}).Info("skip article")
			continue
		}

		post := crawlPost(relation.Content.Id)
		article := generateArticle(post)

		var mediaContent *MediaContent
		for _, mediaRelation := range post.Relations {
			mediaType := mediaRelation.PrimaryType
			if strings.HasSuffix(mediaType, "video") || strings.HasSuffix(mediaType, "audio") {
				mediaContent = &mediaRelation.Content
				break
			}
		}

		if mediaContent == nil {
			continue
		}

		mediaId := mediaContent.ExternalId
		article.MediaId = mediaId

		media := crawlMedia(mediaId)

		var externalMedia *ExternalMedia
		for _, external := range media.Media {
			tmp := external
			if externalMedia == nil {
				externalMedia = &tmp
			} else if tmp.Kind == "video" {
				break
			} else if tmp.Kind == "audio" {
				bitrate, err := strconv.Atoi(tmp.Bitrate)
				if err != nil {
					log.WithFields(logrus.Fields{"bitrate": tmp.Bitrate, "err": err}).Error("bitrate error")
					break
				}
				currentBitrate, err := strconv.Atoi(externalMedia.Bitrate)
				if err != nil {
					log.WithFields(logrus.Fields{"bitrate": externalMedia.Bitrate, "err": err}).Error("bitrate error")
					break
				}
				if bitrate > currentBitrate {
					externalMedia = &tmp
				}
			}
		}

		if externalMedia == nil {
			continue
		}
		article.Type = externalMedia.Kind

		var mediaConn *MediaConnection
		for _, connection := range externalMedia.Connection {
			if connection.Protocol == "http" {
				mediaConn = &connection
				break
			}
		}

		if mediaConn == nil {
			continue
		}

		format := "mp4"
		formatTmp := strings.Split(externalMedia.Type, "/")
		if len(formatTmp) == 2 {
			format = formatTmp[1]
		}

		fileName := mediaContent.ExternalId + "." + format
		filePath := bbcConfig.DownloadPath + fileName
		article.FileName = fileName
		article.FilePath = filePath

		executeCount++
		go downloadAndUpload(article, mediaConn.Href)
	}

	for i := 0; i < executeCount; i++ {
		article := <-uploadChannel
		insertDb(db, article)
	}

	log.Info("crawler done!")
}
