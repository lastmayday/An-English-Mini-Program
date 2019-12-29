package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLog(logFilePath string) *logrus.Logger {
	log := logrus.New()

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.Out = logFile

	return log
}
