package dto

type JsApiTicketOutPut struct {
	// 微信JsSDK签名
	Signature string `json:"signature"`
}

type TokenOutput struct {
	Token string `json:"token"`
}

type GetEnterUrlOutput struct {
	Url string `json:"url"`
}
