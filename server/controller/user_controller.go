package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
)

// users 路由注册
func UserRegister(router *gin.RouterGroup) {
	var (
		userModel      model.UserModel     = model.NewUserModel()
		addrModel      model.UserAddrModel = model.NewUserAddrModel()
		regionModel    model.RegionModel   = model.NewRegionModel()
		userService    service.UserService = service.NewUserService(userModel, addrModel, regionModel)
		userController UserController      = NewUserController(userService)
	)
	router.GET("/profile", middleware.MobileBoundRequired(), userController.UserProfile)

	router.GET("/addrs", middleware.MobileBoundRequired(), userController.ListUserAddr)
	router.POST("/addrs", middleware.MobileBoundRequired(), userController.CreateUserAddr)
	router.GET("/addrs/:id", middleware.MobileBoundRequired(), userController.GetUserAddr)
	router.DELETE("/addrs/:id", middleware.MobileBoundRequired(), userController.DelUserAddr)
	router.PUT("/addrs/:id", middleware.MobileBoundRequired(), userController.UpdateUserAddr)
}

type UserController interface {
	UserProfile(ctx *gin.Context)
	ListUserAddr(ctx *gin.Context)
	CreateUserAddr(ctx *gin.Context)
	GetUserAddr(ctx *gin.Context)
	DelUserAddr(ctx *gin.Context)
	UpdateUserAddr(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

// GetUserDetail godoc
// @Summary 个人中心->账户信息
// @Description 获取用户的profile
// @Tags users
// @Accept json
// @Produce json
// @Param  token header string true "用户token"
// @Success 200 {object} middleware.Response{data=dto.UserDetailOutput}
// @Router /users/profile [get]
func (c *userController) UserProfile(ctx *gin.Context) {
	id := ctx.GetInt64("userId")
	userDetail, err := c.service.Retrieve(id)
	if _, ok := errors.Cause(err).(ecode.Codes); ok {
		util.Log.WithFields(logrus.Fields{
			"user_id": id,
		}).Errorf("获取用户信息出错: err: %s\n", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, userDetail)
	}
}

// ListUserAddr godoc
// @Summary 个人中心->账户信息->收件地址
// @Description 获取收件地址列表
// @Tags addrs
// @Produce json
// @Param  token header string true "用户token"
// @Success 200 {object} middleware.Response{data=[]model.UserAddr}
// @Router /users/addrs [get]
func (c *userController) ListUserAddr(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	addrs, err := c.service.FindAllAddrs(userId)
	if err != nil {
		util.Log.Errorf("获取用户收件地址列表出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, addrs)
	}
}

// CreateUserAddr godoc
// @Summary 新增用户收件地址
// @Description 新增地址
// @Tags addrs
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.CreateUserAddrInput true "新增用户收件地址"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/addrs [post]
func (c *userController) CreateUserAddr(ctx *gin.Context) {
	var addr model.UserAddr
	err := ctx.ShouldBindJSON(&addr)
	if err != nil {
		util.Log.Errorf("参数绑定错误, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	addr.UserId = ctx.GetInt64("userId")

	id, err := c.service.SaveAddr(&addr)
	if err != nil {
		util.Log.Errorf("创建用户收件地址失败, 参数: [%v], err: [%s]", addr, err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
	}
}

// GetUserAddr godoc
// @Summary 获取单个收件地址的详情
// @Description 获取单个收件地址的详情
// @Tags addrs
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "addr id"
// @Success 200 {object} middleware.Response{data=dto.GetUserAddrOutput} "success"
// @Router /users/addrs/{id} [get]
func (c *userController) GetUserAddr(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	addr, err := c.service.RetrieveAddr(id)
	if err != nil {
		util.Log.Errorf("根据id获取addr失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, addr)
	}
}

// GetUserAddr godoc
// @Summary 删除单个收件地址的
// @Description 删除单个收件地址的
// @Tags addrs
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "addr id"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/addrs/{id} [delete]
func (c *userController) DelUserAddr(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	err = c.service.DeleteAddr(id)
	if err != nil {
		util.Log.Errorf("根据id删除addr失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
	}
}

// PutUserAddr godoc
// @Summary 修改单个收件地址的
// @Description 修改单个收件地址的
// @Tags addrs
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "addr id"
// @Param body body dto.UpdateUserAddrInput true "修改用户收件地址"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/addrs/{id} [put]
func (c *userController) UpdateUserAddr(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	var addr dto.UpdateUserAddrInput
	err = ctx.ShouldBindJSON(&addr)
	if err != nil {
		util.Log.Errorf("参数绑定错误, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	err = c.service.UpdateUserAddr(ctx, id, &addr)
	if err != nil {
		util.Log.Errorf("修改用户收件地址失败, id: [%d] 参数: [%v], err: [%s]", id, addr, err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
	}
}

func NewUserController(service service.UserService) UserController {
	return &userController{
		service: service,
	}
}
