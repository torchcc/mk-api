package service

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"mk-api/library/util/cos"
	"mk-api/library/util/sms"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
)

type LoginRegisterService interface {
	GenerateCaptcha(ctx *gin.Context) (captchaImgUrl string, err error)
	GenerateSmsVerificationCode(mobile string) (err error)
	LoginRegister(ctx *gin.Context, input *dto.LoginRegisterInput) (token string, err error)
}

type loginRegisterService struct {
	captchaModel model.CaptchaModel
	userModel    model.UserModel
}

func (service *loginRegisterService) LoginRegister(ctx *gin.Context, input *dto.LoginRegisterInput) (token string, err error) {
	userId := ctx.GetInt64("userId")
	captchaKey := fmt.Sprintf("string.login_captcha.%d", userId)
	smsKey := "string.login_sms." + input.Mobile
	if !service.captchaModel.Check(captchaKey, input.CaptchaCode) {
		return "", errors.New("图形验证码出错")
	}
	if !service.captchaModel.Check(smsKey, input.SmsCode) {
		return "", errors.New("短信验证码出错")
	}
	// 在mysql设置手机号码, 注册经纬度
	if err = service.userModel.AddRegisterInfo(input, userId); err != nil {
		util.Log.Errorf("注册更新手机号码出错, userId: [%d], err: [%s]", userId, err)
		return "", errors.New("服务器内部错误, 请重试")
	}

	openId, err := service.userModel.GetOpenIdByUserId(userId)
	if err != nil {
		util.Log.Errorf("查询open_id出错， user_id: [%d], err: [%s]", userId, err.Error())
		return "", errors.New("服务器内部错误")
	}

	// 更新open_id 对应的userInfo 更新token
	token = service.userModel.UpdateRedisToken(openId, userId, input.Mobile)
	return
}

func (service *loginRegisterService) GenerateCaptcha(ctx *gin.Context) (captchaImgUrl string, err error) {
	img, captchaCode, err := util.GenerateCaptcha()
	if err != nil {
		util.Log.Errorf("生成captcha图片出错, err: [%s]", err.Error())
		return "", err
	}
	// key格式"string.captcha.123"
	key := fmt.Sprintf("string.login_captcha.%d", ctx.GetInt64("userId"))

	go func() {
		err = service.captchaModel.Save(key, captchaCode)
		if err != nil {
			util.Log.Errorf("保存到redis出错, err: [%s]", err.Error())
		}
	}()

	// 将图片上传到cos
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	reader := bytes.NewReader(buf.Bytes())

	captchaImgUrl, err = cos.UploadIOStream("captcha.png", reader, true)

	if err != nil {
		util.Log.Errorf("上传captchaImg到cos出错, err: [%s]", err.Error())
	}

	return
}

func (service *loginRegisterService) GenerateSmsVerificationCode(mobile string) (err error) {
	smsVerificationCode := service.getRandomDigits()
	key := "string.login_sms." + mobile

	// 短信验证码保存在redis
	go func() {
		err = service.captchaModel.Save(key, smsVerificationCode)
		if err != nil {
			util.Log.Errorf("sms保存到redis出错, mobile: [%s], err: [%s]", mobile, err.Error())
			return
		}
	}()

	// 腾讯云发送短信到手机
	go func() {
		err = sms.SendRegisterMsg(mobile, smsVerificationCode)
		if err != nil {
			util.Log.Errorf("腾讯云sms服务出错, mobile: [%s], err: [%s]", mobile, err)
		}
	}()

	return
}

func (service *loginRegisterService) getRandomDigits() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

func NewLoginRegisterService(captchaModel model.CaptchaModel, userModel model.UserModel) LoginRegisterService {
	return &loginRegisterService{
		captchaModel: captchaModel,
		userModel:    userModel,
	}
}
