package logger

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var timestampFormat = time.RFC3339Nano

func SetText() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp: true,
	}
	Logger.Out = os.Stderr
}

func SetJSON() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: timestampFormat,
	}
	Logger.Out = os.Stderr
	Logger.ReportCaller = true
}

func SetQuiet() {
	Logger.Out = ioutil.Discard
}

func SetDebug() {
	Logger.Level = logrus.DebugLevel
}

func init() {
	SetText()
}

func Log() *logrus.Logger {
	return Logger
}
