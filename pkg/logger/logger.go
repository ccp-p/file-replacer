package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Log.SetLevel(logrus.InfoLevel)
}

// SetDebug 设置日志级别为Debug
func SetDebug(debug bool) {
	if debug {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
}
