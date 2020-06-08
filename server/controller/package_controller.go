package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"mk-api/library/ecode"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
)

// packages 路由注册
func PackageRegister(router *gin.RouterGroup) {
	var (
		packageModel      model.PackageModel     = model.NewPackageModel()
		packageService    service.PackageService = service.NewPackageService(packageModel)
		packageController PackageController      = NewPackageController(packageService)
	)
	router.GET("/", packageController.ListPackage)
	router.GET("/:id", packageController.GetPackage)
}

type PackageController interface {
	ListPackage(ctx *gin.Context)
	GetPackage(ctx *gin.Context)
}

type packageController struct {
	service service.PackageService
}

// GetPackageDetail godoc
// @Summary 获取单个套餐详情
// @Description 获取单个套餐详情
// @Tags packages
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "package id"
// @Success 200 {object} middleware.Response{data=dto.GetPackageOutPut} "success"
// @Router /pkg/{id} [get]
func (c *packageController) GetPackage(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	pkg, err := c.service.RetrievePackage(id)
	if err != nil {
		util.Log.Errorf("根据id获取package失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, pkg)
	}
}

// PackageList godoc
// @Summary 获取套餐列表
// @Description 获取套餐列表
// @Tags packages
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param page_size query int false "每页多少条"
// @Param page_no query int false "页码"
// @Param level query string false "医院等级， 0-不限 1-公立三甲 2-公立医院 3-民营医院 4-专业机构"
// @Param category_id query int false "套餐类别id"
// @Param min_price query int false "价格区间左值, 0表示最小"
// @Param max_price query int false "价格区间左值, 0表示无上限"
// @Param target query int false "适用人群 0-不限 1-男士 2-女未婚 3-女已婚"
// @Param disease query int false "高发疾病 0-不限 1-食物不耐受检测，2-骨关节疾病体检 3-健康防癌体检 4-幽门螺旋杆菌检测 5-甲状腺检测 6-糖尿病检测"
// @Param order_by query int false "优先排序 0-默认排序，1-低价优先 2 高价优先"
// @Success 200 {object} middleware.Response{data=dto.PaginateListOutput{list=[]dto.ListPackageOutputEle}}
// @Router /pkg/ [get]
func (c *packageController) ListPackage(ctx *gin.Context) {
	var input dto.ListPackageInput
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		util.Log.Errorf("绑定参数错误 err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	data, err := c.service.ListPackage(ctx, &input)
	if err != nil {
		util.Log.Errorf("获取套餐列表出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, data)
	}
}

func NewPackageController(service service.PackageService) PackageController {
	return &packageController{
		service: service,
	}
}
