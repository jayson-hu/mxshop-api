package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jayson-hu/mxshop-api/user-web/global"
	"github.com/jayson-hu/mxshop-api/user-web/proto"
)

func InitSrvConn()  {
	//从注册中心获取到用户到信息，包括ip和port
	cfg := api.DefaultConfig()
	//cfg.Address = "150.158.11.116:8500"
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0


	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("global.ServerConfig.UserSrvInfo.Name", global.ServerConfig.UserSrvInfo.Name)
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,global.ServerConfig.UserSrvInfo.Name))
	//data, err := client.Agent().ServicesWithFilter(`Service == "user-srv"`)
	if err != nil {

		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	zap.S().Fatal("【init conn 】用户服务不可用")
	fmt.Println("userSrvHost,userSrvPort",userSrvHost,userSrvPort)
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d","127.0.0.1",userSrvPort), grpc.WithInsecure())
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
	//	global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	}
	//存在的问题: 1.后续用户服务下线 2. 改端口 3. 改ip
	//2 多个连接使用同一个连接
	global.UserSrvClient = proto.NewUserClient(userConn)
	
}
