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

func CargoRegister(router *gin.RouterGroup) {
	var (
		cargoModel      model.CargoModel     = model.NewCargoModel()
		packageModel    model.PackageModel   = model.NewPackageModel()
		cargoService    service.CargoService = service.NewCargoService(cargoModel, packageModel)
		cargoController CargoController      = NewCargoController(cargoService)
	)
	router.GET("/", cargoController.GetCargo)
	router.POST("/", cargoController.PostOnePkg2Cargo)
	router.DELETE("/", cargoController.DeleteCargoEntriesByIds)
}

type CargoController interface {
	GetCargo(ctx *gin.Context)
	PostOnePkg2Cargo(ctx *gin.Context)
	DeleteCargoEntriesByIds(ctx *gin.Context)
}

type cargoController struct {
	service service.CargoService
}

// DelCargoEntries godoc
// @Summary 删除购物车条目
// @Description 删除购物车条目
// @Tags cargo
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param  cargo_ids body dto.DeleteCargoEntriesInput true "加购的套餐id"
// @Success 200 {object} middleware.Response{data=string}
// @Router /cargo/ [delete]
func (c *cargoController) DeleteCargoEntriesByIds(ctx *gin.Context) {
	var input dto.DeleteCargoEntriesInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.Log.Errorf("获取cargo_ids参数出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("参数cargo_ids出错"))
		return
	}
	if err := c.service.RemoveCargoEntries(&input); err != nil {
		util.Log.Errorf("删除购物车条目出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("删除购物车条目出错"))
	} else {
		middleware.ResponseSuccess(ctx, "成功")
	}

}

// PostOnePkg2Cargo godoc
// @Summary 往购物车增添一个套餐
// @Description 加购物车
// @Tags cargo
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param  pkgId body dto.PostCargoInput true "加购的套餐id"
// @Success 200 {object} middleware.Response{data=string}
// @Router /cargo/ [post]
func (c *cargoController) PostOnePkg2Cargo(ctx *gin.Context) {
	var input dto.PostCargoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.Log.Errorf("获取pkg_id参数出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("套餐id出错"))
		return
	}
	if err := c.service.CreateCargo(ctx, input.PkgId); err != nil {
		util.Log.Errorf("加购物车出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("加购物车出错"))
	} else {
		middleware.ResponseSuccess(ctx, "成功")
	}
}

// GetCargo godoc
// @Summary 获取购物车里面的pkg列表
// @Description 获取购物车里面的pkg列表
// @Tags cargo
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Success 200 {object} middleware.Response{data=[]dto.GetCargoOutputElem}
// @Router /cargo/ [get]
func (c *cargoController) GetCargo(ctx *gin.Context) {
	pkgs, err := c.service.RetrieveCargo(ctx)
	if err != nil {
		util.Log.Errorf("获取购物车详情失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, pkgs)
	}
}

func NewCargoController(service service.CargoService) CargoController {
	return &cargoController{
		service: service,
	}
}
