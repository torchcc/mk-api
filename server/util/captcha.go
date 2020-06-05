package util

import (
	"image/color"

	"github.com/afocus/captcha"
	"mk-api/server/static"
)

func GenerateCaptcha() (img *captcha.Image, captchaCode string, err error) {
	capGenerator := captcha.New()
	// 设置字体
	err = capGenerator.SetFont(static.Path("font/UniTortred.ttf"))
	if err != nil {
		Log.Errorf("设置字体出错, err: [%s]", err.Error())
		return
	}
	capGenerator.SetSize(128, 64)
	capGenerator.SetDisturbance(captcha.NORMAL)
	capGenerator.SetFrontColor(color.RGBA{0, 0, 0, 255})
	// capGenerator.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	// capGenerator.SetBkgColor(color.RGBA{220, 220, 220, 255})
	// capGenerator.SetBkgColor(color.RGBA{174, 238, 238, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	capGenerator.SetBkgColor(color.RGBA{174, 238, 238, 255})

	// 创建验证码 4个字符 captcha.NUM 字符模式数字类型
	// 返回验证码图像对象以及验证码字符串 后期可以对字符串进行对比 判断验证
	img, captchaCode = capGenerator.Create(4, captcha.NUM)
	return
}
