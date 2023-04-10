package forms

type SendSmsForm struct {
	Mobile    string `json:"mobile" form:"mobile" binding:"required,mobile"` // mobile 与validator指定了，需要在main函数内部进行注册
	Type  string `json:"type" form:"type" binding:"required,oneof=register login"`
	//1.注册发送短信验证码和动态验证码发送验证码
}

