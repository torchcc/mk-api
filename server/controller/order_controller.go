package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/pay"
	payConfig "github.com/silenceper/wechat/v2/pay/config"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/conf"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
)

func OrderRegister(router *gin.RouterGroup) {
	var (
		cartModel    model.CartModel    = model.NewCartModel()
		packageModel model.PackageModel = model.NewPackageModel()
		orderModel   model.OrderModel   = model.NewOrderModel()
		payModel     model.PayModel     = model.NewPayModel()
		cfg                             = &payConfig.Config{
			AppID:     conf.C.WeChat.AppID,
			MchID:     conf.C.WeChat.PayMchID,
			Key:       conf.C.WeChat.PayKey,
			NotifyURL: conf.C.WeChat.PayNotifyURL,
		}
		wechatPay                            = pay.NewPay(cfg)
		orderService    service.OrderService = service.NewOrderService(orderModel, packageModel, cartModel, payModel, wechatPay)
		orderController OrderController      = NewOrderController(orderService)
	)
	router.POST("/orders/", orderController.PostOrder)
	router.GET("/orders/", orderController.ListOrder)
	router.GET("/orders/:id", orderController.GetOrder)
	router.DELETE("/orders/:id", orderController.DeleteOrder)

	router.PUT("/order_items/", orderController.PutOrderItem)
}

type OrderController interface {
	PostOrder(ctx *gin.Context)
	ListOrder(ctx *gin.Context)
	GetOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)

	PutOrderItem(ctx *gin.Context)
}

type orderController struct {
	service service.OrderService
}

// UpdateOrderItem godoc
// @Summary 更新orderItem的体检人信息
// @Description 更新orderItem的体检人信息
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.PutOrderItemInput true "修改的orderItem的体检人信息"
// @Success 200 {object} middleware.Response{data=dto.ResourceID}
// @Router /order_items/ [put]
func (c *orderController) PutOrderItem(ctx *gin.Context) {
	var input dto.PutOrderItemInput
	err := util.ParseRequest(ctx, &input)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	err = c.service.ModifyOrderItem(ctx, &input)
	if err != nil {
		if _, ok := err.(ecode.Codes); ok {
			middleware.ResponseError(ctx, ecode.RequestErr, ctx.Errors.Last())
			util.Log.Errorf("failed to update order item, input is [%v], err is [%c]", input, ctx.Errors.Last())
			return
		}
		util.Log.Errorf("failed to update order item, input is [%v], err is [%c]", input, err)
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("internal server error"))
		return
	}
	middleware.ResponseSuccess(ctx, dto.ResourceID{Id: input.Id})
}

// DeleteOrder godoc
// @Summary 删除订单
// @Description 删除订单
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param id path int true "订单的id, order_id"
// @Success 200 {object} middleware.Response{data=string}
// @Router /orders/{id} [delete]
func (c *orderController) DeleteOrder(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("参数id有误"))
		return
	}
	err = c.service.RemoveOrder(ctx, id)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"order_id": id}).Errorf("移除订单失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, "")
	}
}

// GetOrderDetail godoc
// @Summary 获取订单详情
// @Description 获取订单详情
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param id path int true "订单的id, order_id"
// @Success 200 {object} middleware.Response{data=dto.RetrieveOrderOutput} "success"
// @Router /orders/{id} [get]
func (c *orderController) GetOrder(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("参数id有误"))
		return
	}
	order, err := c.service.RetrieveOrder(ctx, id)
	if err != nil {
		util.Log.Errorf("根据id获取order失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, order)
	}
}

// OrderList godoc
// @Summary 获取订单列表
// @Description 获取订单列表
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param page_size query int false "每页多少条"
// @Param page_no query int false "页码"
// @Param status query int false "订单状态 -1 全部(默认值) 0-未付款，2-已付款(待预约), 3-已退款, 4-已关闭"
// @Success 200 {object} middleware.Response{data=dto.PaginateListOutput{list=[]dto.ListOrderOutputEle}}
// @Router /orders/ [get]
func (c *orderController) ListOrder(ctx *gin.Context) {
	var input dto.ListOrderInput
	if err := util.ParseRequest(ctx, &input); err != nil {
		util.Log.Errorf("参数绑定失败, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	data, err := c.service.ListOrder(ctx, &input)
	if err != nil {
		util.Log.Errorf("获取订单列表失败, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
		return
	}
	middleware.ResponseSuccess(ctx, data)
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 创建订单,返回前端调起微信支付的必须参数
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.PostOrderInput true "创建订单的请求体"
// @Success 200 {object} middleware.Response{data=dto.PostOrderOutput}
// @Router /orders/ [post]
func (c *orderController) PostOrder(ctx *gin.Context) {
	var input dto.PostOrderInput
	err := util.ParseRequest(ctx, &input)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	cfg, err := c.service.CreateOrder(ctx, &input)
	if err != nil {
		util.Log.Errorf("controller failed to create order, err: [%s]", err.Error())
		switch err.(type) {
		case ecode.Codes:
			if ecode.Equal(err.(ecode.Codes), ecode.RequestErr) {
				middleware.ResponseError(ctx, ecode.RequestErr, ctx.Errors.Last())
				return
			}
		}
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	middleware.ResponseSuccess(ctx, cfg)
}

func NewOrderController(service service.OrderService) OrderController {
	return &orderController{
		service: service,
	}
}
