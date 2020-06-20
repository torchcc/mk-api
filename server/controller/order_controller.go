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

func OrderRegister(router *gin.RouterGroup) {
	var (
		cartModel       model.CartModel      = model.NewCartModel()
		packageModel    model.PackageModel   = model.NewPackageModel()
		orderModel      model.OrderModel     = model.NewOrderModel()
		orderService    service.OrderService = service.NewOrderService(orderModel, packageModel, cartModel)
		orderController OrderController      = NewOrderController(orderService)
	)
	router.POST("/", orderController.PostOrder)
}

type OrderController interface {
	PostOrder(ctx *gin.Context)
}

type orderController struct {
	service service.OrderService
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 创建订单
// @Tags orders
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.PostOrderInput true "创建订单的请求体"
// @Success 200 {object} middleware.Response{data=dto.PostOrderOutput} "success"
// @Router /orders/ [post]
func (c *orderController) PostOrder(ctx *gin.Context) {
	var input dto.PostOrderInput
	err := util.ParseRequest(ctx, &input)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	err = c.service.CreateOrder(ctx, &input)
	if err != nil {
		switch err.(type) {
		case ecode.Codes:
			if ecode.Equal(err.(ecode.Codes), ecode.RequestErr) {
				middleware.ResponseError(ctx, ecode.RequestErr, ctx.Errors.Last())
				return
			}
		}
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器错误"))
		return
	}
	// TODO 到这儿了
	ctx.JSON(200, "ok")
}

func NewOrderController(service service.OrderService) OrderController {
	return &orderController{
		service: service,
	}
}
