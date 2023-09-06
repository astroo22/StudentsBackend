package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {

	appEnv := os.Getenv("APP_ENV")

	if appEnv == "prod" {
		// Log to a file in production
		file, err := os.OpenFile("/var/log/backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Log.Fatalf("Failed to open log file: %v", err)
		}
		Log.SetOutput(file)
	} else {
		// Log to standard output (terminal) in other environments
		Log.SetOutput(os.Stdout)
	}
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{})
}
