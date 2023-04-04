package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	url := "https://baidu.com"
	sugar.Infow("infow", "url", url, "attempt", 3)
	sugar.Infof("failed to fetch url: %s", url)
	println("**********")
	logger.Info("logger")

}
