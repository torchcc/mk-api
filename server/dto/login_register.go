package dto

type LoginRegisterInput struct {
	// 用户手机号码
	Mobile string `json:"mobile" description:"用户手机号码" comment:"手机号码" validate:"required"`
	// 图形验证码
	CaptchaCode string `json:"captcha_code" description:"图形验证码" comment:"图形验证码" en_comment:"CaptchaCode" validate:"required"`
	// 短信验证码
	SmsCode string `json:"sms_code" description:"短信验证码" comment:"短信验证码" validate:"required"`
	// 注册时的经度, 获取不到的话请传 0.0
	Longitude float64 `json:"longitude"`
	// 注册时的纬度， 获取不到请传 0.0
	Latitude float64 `json:"latitude"`
}
