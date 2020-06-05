package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"mk-api/library/ecode"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
)

/*登录注册控制器， 未注册的手机号默认注册*/

// login 路由注册
func LoginRegister(router *gin.RouterGroup) {
	var (
		captchaModel            model.CaptchaModel           = model.NewCaptchaModel()
		userModel               model.UserModel              = model.NewUserModel()
		loginRegisterService    service.LoginRegisterService = service.NewLoginRegisterService(captchaModel, userModel)
		loginRegisterController LoginRegisterController      = NewLoginRegisterController(loginRegisterService)
	)
	router.GET("/captcha", loginRegisterController.GetCaptchaImg)
	router.GET("/sms", loginRegisterController.GetSmsVerificationCode)
	router.POST("/", loginRegisterController.LoginOrRegister)
}

type LoginRegisterController interface {
	LoginOrRegister(ctx *gin.Context)
	GetCaptchaImg(ctx *gin.Context)
	GetSmsVerificationCode(ctx *gin.Context)
}

type loginRegisterController struct {
	service service.LoginRegisterService
}

// LoginRegister godoc
// @Summary 登录或者注册
// @Description 登录或者注册
// @Tags login
// @Accept json
// @Produce json
// @Param token header string true "用户token"
// @Param loginBody body dto.LoginRegisterInput true "captcha_code:图形验证码， sms_code:短信验证码"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /login_register/ [post]
func (c *loginRegisterController) LoginOrRegister(ctx *gin.Context) {
	// 该接口已经测试
	var loginPayload dto.LoginRegisterInput
	err := ctx.ShouldBindJSON(&loginPayload)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	err = c.service.LoginRegister(ctx, &loginPayload)
	if err != nil {
		util.Log.Errorf("登录/注册出错, 参数：[%v], err: [%s]", loginPayload, err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, nil)
	}
}

// GetCaptcha godoc
// @Summary 获取图片验证码
// @Description 获取图片验证码
// @Tags login
// @Accept json
// @Produce json
// @Param token header string true "用户token"
// @Success 200 {object} middleware.Response{data=dto.GetCaptchaOutput}
// @Router /login_register/captcha [get]
func (c *loginRegisterController) GetCaptchaImg(ctx *gin.Context) {
	captchaImgUrl, err := c.service.GenerateCaptcha(ctx)
	if err != nil {
		util.Log.Errorf("生成captcha失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("生成验证码失败，请重试"))
	} else {
		middleware.ResponseSuccess(ctx, dto.GetCaptchaOutput{CaptchaImgUrl: captchaImgUrl})
	}
}

// GetSMS godoc
// @Summary 获取短信验证码
// @Description 获取短信验证码
// @Tags login
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param mobile query string true "用户手机号码"
// @Success 200 {object} middleware.Response{}
// @Router /login_register/sms [get]
func (c *loginRegisterController) GetSmsVerificationCode(ctx *gin.Context) {
	mobile := ctx.Query("mobile")
	if !util.IsMobile(mobile) {
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("手机号码格式不正确"))
		return
	}

	err := c.service.GenerateSmsVerificationCode(mobile)
	if err != nil {
		util.Log.Errorf("发送短信验证码错误：mobile: [%s], err: [%s]", mobile, err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, nil)
	}
}

func NewLoginRegisterController(service service.LoginRegisterService) LoginRegisterController {
	return &loginRegisterController{
		service: service,
	}
}
