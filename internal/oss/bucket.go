package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

type OssConfig struct {
	EndPoint        string `mapstructure:"ENDPOINT"`
	AccessKeyId     string `mapstructure:"ACCESS_KEY_ID"`
	AccessKeySecret string `mapstructure:"ACCESS_KEY_SECRET"`
	BucketName      string `mapstructure:"BUCKET_NAME"`
	BaseOssUrl      string `mapstructure:"BASE_OSS_URL"`
}

var ossConfig OssConfig

func init() {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("oss_config")
	v.AddConfigPath("$HOME/go/src/github.com/lastmayday/BBC-English/config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&ossConfig)
	if err != nil {
		panic(err)
	}
}

func Upload(fileName string, filePath string) (string, error) {
	client, err := oss.New(ossConfig.EndPoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(ossConfig.BucketName)
	if err != nil {
		return "", err
	}
	err = bucket.PutObjectFromFile(fileName, filePath)
	if err != nil {
		return "", err
	}

	ossUrl := ossConfig.BaseOssUrl + fileName
	return ossUrl, nil
}
