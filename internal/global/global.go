package global

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	Log = logger.Sugar()
}
