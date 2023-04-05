package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//使用正则表达是否合法
	ok, _ := regexp.MatchString(`^1\d{10}$`, mobile)
	//ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		//fmt.Println("ok", ok)
		return false
	}
	//fmt.Println("17ok", ok)
	return true

}
