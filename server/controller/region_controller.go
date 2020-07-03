package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"mk-api/library/ecode"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
)

func RegionRegister(router *gin.RouterGroup) {
	var (
		regionModel      model.RegionModel     = model.NewRegionModel()
		regionService    service.RegionService = service.NewRegionService(regionModel)
		regionController RegionController      = NewRegionController(regionService)
	)
	router.GET("/", regionController.ListRegion)
}

type RegionController interface {
	ListRegion(ctx *gin.Context)
}

type regionController struct {
	service service.RegionService
}

// ListRegion godoc
// @Summary 根据parent_id获取行政区域列表
// @Description 根据parent_id获取行政区域列表
// @Tags regions
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Param  parent_id query int false "父级行政区域id， 默认值为第一级(省等)：0"
// @Success 200 {object} middleware.Response{data=[]dto.Region}
// @Router /regions/ [get]
func (c *regionController) ListRegion(ctx *gin.Context) {
	parentId, err := strconv.ParseInt(ctx.Query("parent_id"), 10, 64)
	if err != nil {
		parentId = int64(0)
	}
	output, err := c.service.RetrieveRegionsByParentId(parentId)
	if err != nil {
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("内部服务器错误"))
		return
	}
	middleware.ResponseSuccess(ctx, output)

}

func NewRegionController(service service.RegionService) RegionController {
	return &regionController{
		service: service,
	}
}
