package initialize

import (
	"go.uber.org/zap"
)

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("init logger failed, err: " + err.Error())
	}

	zap.ReplaceGlobals(logger)
}
