package utils

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestZapSugarDemo(t *testing.T) {
	//logger, _ := zap.NewDevelopment() // 开发环境
	logger, _ := zap.NewProduction() // 生产环境
	defer logger.Sync()              // flushes buffer, if any

	sugar := logger.Sugar()
	url := "https://imooc.com"
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

func TestZapDemo(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	url := "https://imooc.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stderr",
	}

	return cfg.Build()
}

func TestZapLogFile(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		panic("init logger failed: " + err.Error())
	}

	su := logger.Sugar()
	defer su.Sync()

	url := "https://imooc.com"
	su.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
	)
}
