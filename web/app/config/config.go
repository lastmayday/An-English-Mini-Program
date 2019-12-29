package config

import (
	"github.com/spf13/viper"
)

type AppConfigStruct struct {
	AppLogFile      string `mapstructure:"APP_LOG_FILE"`
	AppPort         string `mapstructure:"APP_PORT"`
	ArticlePageSize int    `mapstructure:"ARTICLE_PAGE_SIZE"`
}

var AppConfig AppConfigStruct

func init() {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("app_config")
	v.AddConfigPath("$HOME/go/src/github.com/lastmayday/BBC-English/config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}
}
