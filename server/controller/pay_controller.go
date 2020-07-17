package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	payConf "github.com/silenceper/wechat/v2/pay/config"
	"github.com/silenceper/wechat/v2/pay/notify"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/conf"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
)

func PayRegister(router *gin.RouterGroup) {
	var (
		payModel   model.PayModel   = model.NewPayModel()
		orderModel model.OrderModel = model.NewOrderModel()
		cfg                         = &payConf.Config{
			AppID:     conf.C.WeChat.AppID,
			MchID:     conf.C.WeChat.PayMchID,
			Key:       conf.C.WeChat.PayKey,
			NotifyURL: conf.C.WeChat.PayNotifyURL,
		}
		ntf           *notify.Notify     = notify.NewNotify(cfg)
		payService    service.PayService = service.NewPayService(ntf, payModel, orderModel)
		payController PayController      = NewPayController(payService)
	)
	router.POST("/wechat_callback", payController.WechatPayCallback)
	router.GET("/status", middleware.MobileBoundRequired(), payController.CheckPayStatus)
	router.POST("/scnd_pay", middleware.MobileBoundRequired(), payController.Launch2ndPay)
}

type PayController interface {
	WechatPayCallback(ctx *gin.Context)
	CheckPayStatus(ctx *gin.Context)
	Launch2ndPay(ctx *gin.Context)
}

type payController struct {
	service service.PayService
}

// CreateOrder godoc
// @Summary 根据订单id发起二次支付
// @Description 根据订单id发起二次支付,返回前端调起微信支付的必须参数
// @Tags pay
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.ResourceID true "创建订单的请求体"
// @Success 200 {object} middleware.Response{data=dto.PostOrderOutput}
// @Router /pay/scnd_pay [post]
func (c *payController) Launch2ndPay(ctx *gin.Context) {
	var order dto.ResourceID
	if err := ctx.ShouldBindJSON(&order); err != nil {
		util.Log.Error("绑定参数错误")
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("参数id出错"))
		return
	}
	cfg, err := c.service.Launch2ndPay(ctx, order.Id)
	if err != nil {
		switch err.(type) {
		case ecode.Codes:
			if ecode.Equal(err.(ecode.Codes), ecode.RequestErr) {
				util.Log.WithFields(logrus.Fields{"order_id": order.Id}).
					Warningf("重新支付失败, err: [%s]", ctx.Errors.Last().Error())
				middleware.ResponseError(ctx, ecode.RequestErr, ctx.Errors.Last())
				return
			}
			if ecode.Equal(err.(ecode.Codes), ecode.NothingFound) {
				util.Log.WithFields(logrus.Fields{"order_id": order.Id}).
					Warningf("重新支付失败, err: [%s]", ctx.Errors.Last().Error())
				middleware.ResponseError(ctx, ecode.NothingFound, errors.New("啥都木有"))
				return
			}
		}
		util.Log.WithFields(logrus.Fields{"order_id": order.Id}).
			Errorf("重新支付失败, err: [%s]", ctx.Errors.Last().Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器错误"))
	}
	middleware.ResponseSuccess(ctx, cfg)
}

// CheckPayStatus godoc
// @Summary 查询订单支付状态
// @Description 前端轮询支付状态
// @Tags pay
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param prepay_id query string true "订单的prepay_id"
// @Success 200 {object} middleware.Response{data=dto.CheckPayStatusOutput}
// @Router /pay/status [get]
func (c *payController) CheckPayStatus(ctx *gin.Context) {
	prepayId := ctx.Query("prepay_id")
	if prepayId == "" {
		errStr := "请求参数prepay_id有误"
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New(errStr))
		util.Log.Warning(errStr)
		return
	}

	status, err := c.service.CheckPayStatus(ctx, prepayId)
	if err != nil {
		util.Log.Errorf("查询支付状态出错，err: [%s]", err)
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("内部服务器错误"))
		return
	}
	middleware.ResponseSuccess(ctx, dto.CheckPayStatusOutput{Status: status})

}

// 微信支付回调 Notify
func (c *payController) WechatPayCallback(ctx *gin.Context) {
	const AckSuccess = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	const AckFail = "<xml><return_code><![CDATA[FAIL]]></return_code></xml>"
	// ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if ok := c.service.WechatPayCallBack(ctx); ok {
		_, _ = fmt.Fprint(ctx.Writer, AckSuccess)
		return
	}
	_, _ = fmt.Fprint(ctx.Writer, AckFail)
}

func NewPayController(service service.PayService) PayController {
	return &payController{
		service: service,
	}
}
