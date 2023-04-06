package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/jayson-hu/mxshop-api/user-web/middlewares"
	user_router "github.com/jayson-hu/mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/v1")

	user_router.InitUserRouter(ApiGroup)
	return Router
}
