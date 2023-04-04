package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jayson-hu/mxshop-api/user-web/api"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	//UserRouter.GET("list", api.GetUserList())
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
