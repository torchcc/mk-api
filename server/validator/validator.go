package validator

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"mk-api/server/validator/id_card"
)

var (
	V     *validator.Validate
	Trans ut.Translator
)

func checkMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobile)
}

func checkIdCardNo(fl validator.FieldLevel) bool {
	idCardNo := []byte(fl.Field().String())
	return id_card.IsValidCitizenNo(&idCardNo)
}

func Init() {
	// 中文翻译
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	Trans, _ = uni.GetTranslator("zh")

	V, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 验证器注册翻译器
		_ = zhTranslations.RegisterDefaultTranslations(V, Trans)

		// 自定义验证方法 验证手机
		{
			_ = V.RegisterValidation("checkMobile", checkMobile)

			_ = V.RegisterTranslation("checkMobile", Trans, func(ut ut.Translator) error {
				return ut.Add("checkMobile", "{0}长度不等于11位或{1}格式错误!", true) // see universal-translator for details
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("checkMobile", fe.Field(), fe.Field())
				return t
			})
		}

		// 验证身份证号码
		{
			_ = V.RegisterValidation("checkIdCardNo", checkIdCardNo)

			_ = V.RegisterTranslation("checkIdCardNo", Trans, func(ut ut.Translator) error {
				return ut.Add("checkIdCardNo", "身份证号码有误", true) // see universal-translator for details
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("checkIdCardNo", fe.Field(), fe.Field())
				return t
			})
		}
	}
}

func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a xtime.
		errList = append(errList, e.Translate(Trans))
	}
	return strings.Join(errList, "|")
}
