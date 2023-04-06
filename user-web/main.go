package main

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/jayson-hu/mxshop-api/user-web/global"
	"github.com/jayson-hu/mxshop-api/user-web/initialize"
	userwebvalidator "github.com/jayson-hu/mxshop-api/user-web/validator"
)

func main() {
	//port := 8022

	//初始化logger
	initialize.InitLogger()
	//初始化配置文件
	initialize.InitConfig()
	port := global.ServerConfig.Port
	//初始化trans
	if err := initialize.InitTranslator("zh"); err != nil {
		panic(err)
	}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", userwebvalidator.ValidateMobile)
		//针对于trans 翻译错误， 未能正常处理英文的错误
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	//初始化router
	router := initialize.Routers()
	zap.S().Debugf("======启动服务，端口为： %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panicf("启动失败： ", err.Error())
	}
}
