package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

const (
	jsonMode    = "LOG_JSON"
	EnvLevel    = "LOG_LEVEL"
	ServiceName = "sow_library"
)

func NewLogger() *Logger {
	inst := logrus.New()

	envLogJson := strings.Trim(strings.ToLower(os.Getenv(jsonMode)), " ")

	logLevel, err := logrus.ParseLevel(os.Getenv(EnvLevel))
	if err != nil {
		inst.Warnf("Logger: %s", err)
		logLevel = logrus.InfoLevel
	}

	if envLogJson == "true" {
		inst.SetFormatter(&logrus.JSONFormatter{})
		inst.Warningln("JSON mode enabled!")
	}

	inst.SetOutput(os.Stdout)
	inst.SetLevel(logLevel)
	inst.WithField("name", ServiceName)
	inst.Warningf("Log EnvLevel: %s", logLevel.String())

	return &Logger{
		Logger: inst,
	}
}
