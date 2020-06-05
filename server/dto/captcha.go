package dto

type GetCaptchaOutput struct {
	CaptchaImgUrl string `json:"captcha_img_url" comment:"captcha image url" description:"captcha image url description"`
}
