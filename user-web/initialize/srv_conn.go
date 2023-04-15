package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jayson-hu/mxshop-api/user-web/global"
	"github.com/jayson-hu/mxshop-api/user-web/proto"
)

// InitSrvConn2 version 2
func InitSrvConn2()  {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		//fmt.Sprintf("consul://%s:%d/%s?wait=14s","150.158.11.116",8500,"user-srv"),
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",consulInfo.Host,consulInfo.Port,global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("srvConn 连接用户服务失败", "msg:", err.Error())
	}
	defer userConn.Close()
	global.UserSrvClient = proto.NewUserClient(userConn)
}

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
	//fmt.Println("global.ServerConfig.UserSrvInfo.Name", global.ServerConfig.UserSrvInfo.Name)
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
	//zap.S().Fatal("【init conn 】用户服务不可用")
	//fmt.Println("userSrvHost,userSrvPort",userSrvHost,userSrvPort)
	zap.S().Infof("获取的userSrvHost: %s(还是使用localhost), userSrvPort: %d",userSrvHost,userSrvPort)
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d","127.0.0.1",userSrvPort), grpc.WithInsecure())
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
	//	global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("SrvConn 连接用户服务失败", "msg:", err.Error())
	}
	//存在的问题: 1.后续用户服务下线 2. 改端口 3. 改ip
	//2 多个连接使用同一个连接
	global.UserSrvClient = proto.NewUserClient(userConn)
	
}
