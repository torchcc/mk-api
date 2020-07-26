package dto

type JsApiTicketOutPut struct {
	// 微信JsSDK签名
	Signature string `json:"signature"`
}

type TokenOutput struct {
	Token string `json:"token"`
	// 是否已经验证了手机号码
	MobileVerified int8 `json:"mobile_verified"`
}

type GetEnterUrlOutput struct {
	Url string `json:"url"`
}
