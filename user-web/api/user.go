package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jayson-hu/mxshop-api/user-web/forms"
	"github.com/jayson-hu/mxshop-api/user-web/global"
	"github.com/jayson-hu/mxshop-api/user-web/global/response"
	"github.com/jayson-hu/mxshop-api/user-web/middlewares"
	"github.com/jayson-hu/mxshop-api/user-web/models"
	"github.com/jayson-hu/mxshop-api/user-web/proto"
)

func RemoveTopStruct(filds map[string]string) map[string]string {
	// User.Name 去除User
	resp := map[string]string{}
	for filed, err := range filds {
		//fmt.Println("333")
		//fmt.Println(filed, err)
		//fmt.Println("ddd333remove",strings.Index(filed, "." ))
		resp[filed[strings.Index(filed, ".")+1:]] = err
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
					//"msg": e.Code(),
					//"msg": e.Message(),
					"msg": e.Code(),
				})
			}
			return
		}
	}

}

func HandleValidatorError(c *gin.Context, err error) {
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
	////从注册中心获取到用户到信息，包括ip和port
	//cfg := api.DefaultConfig()
	////cfg.Address = "150.158.11.116:8500"
	//consulInfo := global.ServerConfig.ConsulInfo
	//cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	//
	//userSrvHost := ""
	//userSrvPort := 0
	//
	//
	//client, err := api.NewClient(cfg)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("global.ServerConfig.UserSrvInfo.Name", global.ServerConfig.UserSrvInfo.Name)
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,global.ServerConfig.UserSrvInfo.Name))
	////data, err := client.Agent().ServicesWithFilter(`Service == "user-srv"`)
	//if err != nil {
	//
	//	panic(err)
	//}
	//for _, value := range data {
	//	userSrvHost = value.Address
	//	userSrvPort = value.Port
	//	break
	//}
	//if userSrvHost == ""{
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"captcha":"用户服务不可大",
	//	})
	//}
	//
	//
	//
	////ip := "127.0.0.1"
	////port := 50051
	//
	////userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())， 拨号连接服务器
	//fmt.Println("userSrvHost,userSrvPort",userSrvHost,userSrvPort)
	////127.0.0.1 为userSrvHost,但是健康检查不可用，会被删除，故使用本地，写死
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d","127.0.0.1",userSrvPort), grpc.WithInsecure())
	////userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
	////	global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	//}
	//生成grpc的client并调用借口
	//userSrvClient := proto.NewUserClient(userConn)
	//以上使用了全局的

	//jwt中claim
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "10")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
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
	//进行验证码的验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 验证成功后的返回
	//拨号连接服务器
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
	//	global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	//}
	//userSrvClient := proto.NewUserClient(userConn)
	//以上使用了全局的srv

	//登录的逻辑
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				zap.S().Infof("登录失败: %s", err.Error())
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败1",
				})
			}
			return
		}
	} else {
		//只是查询到了用户而已， 并没有检查密码
		if passRsp, passErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败2",
			})
			return
		} else {
			//验证密码返回是否为success true 或者是false
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "user-web",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, map[string]string{
						"message": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusOK, map[string]string{
					"message": "密码错误",
				})
			}
		}
	}
}

func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	//存储验证码, 省略了验证短信
	//rdb := redis.NewClient(&redis)

	//拨号连接服务器
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
	//	global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg:", err.Error())
	//}
	//userSrvClient := proto.NewUserClient(userConn)
	//使用了全局的usersrvclient
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		//fmt.Println("msg, 用户失败 ", err.Error())
		zap.S().Errorw("[Register fail ] , 新建用户失败", "msg:", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "user-web",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
