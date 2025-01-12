package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger sets up and returns a configured Logrus logger
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
