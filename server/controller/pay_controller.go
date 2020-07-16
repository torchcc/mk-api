package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	payConf "github.com/silenceper/wechat/v2/pay/config"
	"github.com/silenceper/wechat/v2/pay/notify"
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
}

type PayController interface {
	WechatPayCallback(ctx *gin.Context)
	CheckPayStatus(ctx *gin.Context)
}

type payController struct {
	service service.PayService
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
	const AckSuccess = "<xml><return_code><![CDATA[SUCCESS]]></return_code></xml>"
	const S = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	const AckFail = "<xml><return_code><![CDATA[FAIL]]></return_code></xml>"
	// ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if ok := c.service.WechatPayCallBack(ctx); ok {
		_, _ = fmt.Fprint(ctx.Writer, S)
		return
	}
	_, _ = fmt.Fprint(ctx.Writer, AckFail)
}

func NewPayController(service service.PayService) PayController {
	return &payController{
		service: service,
	}
}
