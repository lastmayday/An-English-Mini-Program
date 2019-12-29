package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Database struct {
	DbUrl      string `mapstructure:"DB_URL"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
}

var DbConfig Database

func init() {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("db_config")
	v.AddConfigPath("$HOME/go/src/github.com/lastmayday/BBC-English/config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&DbConfig)
	if err != nil {
		panic(err)
	}
}

func CreateDbConn() (*sql.DB, error) {
	dbUrl := DbConfig.DbUser + ":" + DbConfig.DbPassword + "@(" + DbConfig.DbUrl + ":" + DbConfig.DbPort + ")/" + DbConfig.DbName
	return sql.Open("mysql", dbUrl)
}
