package service

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitializeLogger() (*logrus.Logger, error) {
	if Log == nil {
		Log = logrus.New()

		logPath := os.Getenv("LOG_PATH")
		logLevel := os.Getenv("LOG_LEVEL")
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			return nil, err
		}

		Log.SetOutput(file)
		logrusLevel, err := logrus.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			return nil, err
		}

		Log.SetLevel(logrusLevel)
	}

	return Log, nil
}
