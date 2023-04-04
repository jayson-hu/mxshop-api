package main

import (
	"fmt"
	"github.com/jayson-hu/mxshop-api/user-web/global"

	"go.uber.org/zap"

	"github.com/jayson-hu/mxshop-api/user-web/initialize"
)

func main() {
	//port := 8022
	//初始化logger
	initialize.InitLogger()
	//初始化配置文件
	initialize.InitConfig()
	port := global.ServerConfig.Port
	//初始化router
	router := initialize.Routers()

	zap.S().Debugf("======启动服务，端口为： %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panicf("启动失败： ", err.Error())
	}
}
