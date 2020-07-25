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

func CartRegister(router *gin.RouterGroup) {
	var (
		cartModel      model.CartModel     = model.NewCartModel()
		packageModel   model.PackageModel  = model.NewPackageModel()
		cartService    service.CartService = service.NewCartService(cartModel, packageModel)
		cartController CartController      = NewCartController(cartService)
	)
	router.GET("/", cartController.GetCart)
	router.POST("/", cartController.PostOnePkg2Cart)
	router.DELETE("/", cartController.DeleteCartEntriesByIds)
}

type CartController interface {
	GetCart(ctx *gin.Context)
	PostOnePkg2Cart(ctx *gin.Context)
	DeleteCartEntriesByIds(ctx *gin.Context)
}

type cartController struct {
	service service.CartService
}

// DelCartEntries godoc
// @Summary 删除购物车条目
// @Description 删除购物车条目
// @Tags cart
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param  cart_ids body dto.DeleteCartEntriesInput true "要删除的cart_id列表"
// @Success 200 {object} middleware.Response{data=string}
// @Router /cart/ [delete]
func (c *cartController) DeleteCartEntriesByIds(ctx *gin.Context) {
	var input dto.DeleteCartEntriesInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.Log.Errorf("获取cart_ids参数出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("参数cart_ids出错"))
		return
	}
	if err := c.service.RemoveCartEntries(&input); err != nil {
		util.Log.Errorf("删除购物车条目出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("删除购物车条目出错"))
	} else {
		middleware.ResponseSuccess(ctx, "成功")
	}

}

// PostOnePkg2Cart godoc
// @Summary 往购物车增添一个套餐
// @Description 加购物车
// @Tags cart
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param  pkgId body dto.PostCartInput true "加购的套餐id"
// @Success 200 {object} middleware.Response{data=string}
// @Router /cart/ [post]
func (c *cartController) PostOnePkg2Cart(ctx *gin.Context) {
	var input dto.PostCartInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.Log.Errorf("获取pkg_id参数出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("套餐id出错"))
		return
	}
	if err := c.service.CreateCart(ctx, input.PkgId); err != nil {
		util.Log.Errorf("加购物车出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("加购物车出错"))
	} else {
		middleware.ResponseSuccess(ctx, "成功")
	}
}

// GetCart godoc
// @Summary 获取购物车里面的pkg列表
// @Description 获取购物车里面的pkg列表
// @Tags cart
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Success 200 {object} middleware.Response{data=[]dto.GetCartOutputElem}
// @Router /cart/ [get]
func (c *cartController) GetCart(ctx *gin.Context) {
	pkgs, err := c.service.RetrieveCart(ctx)
	if err != nil {
		util.Log.Errorf("获取购物车详情失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, pkgs)
	}
}

func NewCartController(service service.CartService) CartController {
	return &cartController{
		service: service,
	}
}
