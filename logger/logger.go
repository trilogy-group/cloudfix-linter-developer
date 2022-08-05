package logger

import (
	"errors"
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var logUser = logrus.New()
var logDev = logrus.New()

func InitLogger(logpath string, jsonFlag bool) error {
	errDir := os.MkdirAll(filepath.Join(logpath, ".cloudfix-linter"), os.ModePerm)
	if errDir != nil {
		return errors.New("can't create cloudfix-linter dir")
	}
	logOutputFile, err := os.OpenFile(filepath.Join(logpath, ".cloudfix-linter", "debug-"+uuid.New().String()+".json"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}
	logDev.SetOutput(logOutputFile)
	logDev.SetReportCaller(true)
	logDev.SetFormatter(&logrus.JSONFormatter{})
	
	if jsonFlag {
		// Turn off the user logger completely in json mode because the only output should be json.
		logUser.SetOutput(ioutil.Discard)
	} else {
		logUser.SetOutput(os.Stdout)
	}
	logUser.SetFormatter(&logrus.TextFormatter{})
	return nil
}

func Info(args ...interface{}) {
	logUser.Info(args)
}

func DevLogger() *logrus.Logger {
	return logDev
}
