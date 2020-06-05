package service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/afocus/captcha"
	"mk-api/server/static"
)

func TestCaptcha(t *testing.T) {
	cap := captcha.New()
	// 设置字体
	pa := static.Path("font/UniTortred.ttf")
	cap.SetFont(pa)
	// 创建验证码 4个字符 captcha.NUM 字符模式数字类型
	// 返回验证码图像对象以及验证码字符串 后期可以对字符串进行对比 判断验证
	img, str := cap.Create(4, captcha.NUM)
	t.Log(img)
	t.Logf("the str is : [%s]", str)

}

func TestGetRandomDigits(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Logf(fmt.Sprintf("%06v", rnd.Int31n(1000000)))
}
