package factory

import (
	"fmt"

	"github.com/trinitytechnology/ebrick/logger/core"
	"github.com/trinitytechnology/ebrick/logger/logrus"
	"github.com/trinitytechnology/ebrick/logger/syslog"
	"github.com/trinitytechnology/ebrick/logger/zap"
)

// LoggerFactory is responsible for creating logger instances based on type and environment.
func LoggerFactory(loggerType, env string, fields ...core.Field) (core.Logger, error) {
	switch loggerType {
	case "zap":
		return zap.NewZapLogger(env, fields...)
	case "logrus":
		return logrus.NewLogrusLogger(env, fields...)
	case "syslog":
		return syslog.NewSysLogger()
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", loggerType)
	}
}
