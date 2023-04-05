package forms

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"` // mobile 与validator指定了，需要在main函数内部进行注册
	Password string `json:"password" form:"password" binding:"required,min=3,max=20"`
}
