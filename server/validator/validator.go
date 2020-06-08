package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var V *validator.Validate

// ValidateCoolTitle returns true when the field value contains the word "cool".
func checkMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobile)
}

func Init() {
	V, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 自定义验证方法
		_ = V.RegisterValidation("checkMobile", checkMobile)
	}
}
