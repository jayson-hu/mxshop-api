package api

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jayson-hu/mxshop-api/user-web/forms"
	"github.com/jayson-hu/mxshop-api/user-web/global"
	"github.com/jayson-hu/mxshop-api/user-web/global/response"
	"github.com/jayson-hu/mxshop-api/user-web/proto"
)

func RemoveTopStruct(filds map[string]string) map[string]string {
	// User.Name 去除User
	resp := map[string]string{}
	for filed, err := range filds {
		//fmt.Println("333")
		//fmt.Println(filed, err)
		//fmt.Println("ddd333remove",strings.Index(filed, "." ))
		resp[filed[strings.Index(filed, "." )+1:]] = err
	}
	return resp
}

// HandleGrpcErrorToHttp 将grpc 的code转换为http的状态吗
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					//"msg":"其他错误" + e.Message(),
					"msg": e.Code(),
				})
			}
			return
		}
	}

}

func HandleValidatorError(c *gin.Context, err error)  {
	errInfo, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	//fmt.Println("errInfo", errInfo)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errInfo.Translate(global.Trans)),
	})
	return

}


func GetUserList(ctx *gin.Context) {
	//ip := "127.0.0.1"
	//port := 50051

	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())， 拨号连接服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "10")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//result := make(map[string]interface{}, 0)
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {

		//data := make(map[string]interface{})
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)),
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02 15:04:05"),
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 验证成功后的返回
	//拨号连接服务器
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)

	//登录的逻辑
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
				"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
				"mobile": "登录失败",
				})
			}
			return
		}
	}else {
		//只是查询到了用户而已， 并没有检查密码
		if passRsp, passErr := userSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password":"登录失败",
			})
			return
		}else {
			//验证密码返回是否为success true 或者是false
			if passRsp.Success {
				c.JSON(http.StatusOK, map[string]string{
					"message":"登录成功",
				})
			}else {
				c.JSON(http.StatusOK, map[string]string{
					"message":"密码错误",
				})
			}
		}
	}
}
