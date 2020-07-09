package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/library/util/cos"
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
		examineeModel  model.ExamineeModel = model.NewExamineeModel()
		userService    service.UserService = service.NewUserService(userModel, addrModel, regionModel, examineeModel)
		userController UserController      = NewUserController(userService)
	)
	router.GET("/profile", userController.GETUserProfile)
	router.PUT("/profile", userController.PutUserProfile)
	router.POST("/profile/avatar", userController.UploadAvatar)

	router.GET("/addrs", userController.ListUserAddr)
	router.POST("/addrs", userController.PostUserAddr)
	router.GET("/addrs/:id", userController.GetUserAddr)
	router.DELETE("/addrs/:id", userController.DelUserAddr)
	router.PUT("/addrs/:id", userController.PutUserAddr)

	router.GET("/examinees", userController.ListExaminee)
	router.POST("/examinees", userController.PostExaminee)
	router.DELETE("/examinees/:id", userController.DelExaminee)
	router.PUT("/examinees/:id", userController.PutExaminee)
}

type UserController interface {
	GETUserProfile(ctx *gin.Context)
	PutUserProfile(ctx *gin.Context)
	UploadAvatar(ctx *gin.Context)

	ListUserAddr(ctx *gin.Context)
	PostUserAddr(ctx *gin.Context)
	GetUserAddr(ctx *gin.Context)
	DelUserAddr(ctx *gin.Context)
	PutUserAddr(ctx *gin.Context)

	ListExaminee(ctx *gin.Context)
	PostExaminee(ctx *gin.Context)
	DelExaminee(ctx *gin.Context)
	PutExaminee(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

// PutExaminee godoc
// @Summary 上传个人头像
// @Description 上传个人头像
// @Tags users
// @accept multipart/form-data
// @Produce  application/json
// @Param token header string true "用户token"
// @Param avatar formData file true "用户上传头像"
// @Success 200 {object} middleware.Response{data=dto.UploadUserAvatarOutput} "success"
// @Router /users/profile/avatar [post]
func (c *userController) UploadAvatar(ctx *gin.Context) {
	_, avatar, err := ctx.Request.FormFile("avatar")
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	err, filePath, _ := cos.Upload2QiNiu(avatar)
	if err != nil {
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("接受返回值失败"))
		return
	}
	err = c.service.UploadAvatar(ctx, filePath)
	if err != nil {
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("修改数据库链接失败"))
		return
	}
	middleware.ResponseSuccess(ctx, dto.UploadUserAvatarOutput{AvatarUrl: filePath})

}

// PutExaminee godoc
// @Summary 修改个人信息
// @Description 修改个人信息
// @Tags users
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.PutUserProfileInput true "修改个人信息"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/profile [put]
func (c *userController) PutUserProfile(ctx *gin.Context) {
	var input dto.PutUserProfileInput
	if err := util.ParseRequest(ctx, &input); err != nil {
		util.Log.Errorf("参数绑定失败, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	if err := c.service.ModifyProfile(ctx, &input); err != nil {
		util.Log.Errorf("修改用户信息失败， payload: [%v]", input)
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("内部服务器出错"))
		return
	}
	middleware.ResponseSuccess(ctx, dto.ResourceID{Id: input.UserId})
}

// PutExaminee godoc
// @Summary 修改单个常用体检人
// @Description 修改单个常用体检人
// @Tags examinees
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "examinee id"
// @Param body body dto.PostExamineeInput true "修改单个常用体检人"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/examinees/{id} [put]
func (c *userController) PutExaminee(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	var input dto.PostExamineeInput

	if err := util.ParseRequest(ctx, &input); err != nil {
		util.Log.Errorf("参数绑定失败, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	err = c.service.ModifyExaminee(ctx, id, &input)
	if err != nil {
		util.Log.Errorf("修改用户收件地址失败, id: [%d] 参数: [%v], err: [%s]", id, input, err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
	}
}

// RemoveExaminee godoc
// @Summary 删除常用体检人信息
// @Description 删除常用体检人信息
// @Tags examinees
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param  id path int true "examinee id"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/examinees/{id} [delete]
func (c *userController) DelExaminee(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	userId := ctx.GetInt64("userId")
	err = c.service.RemoveExaminee(id, userId)
	if err != nil {
		util.Log.Errorf("根据id删除examinee失败, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
	}
}

// PostExaminee godoc
// @Summary 新增常用体检人
// @Description 新增常用体检人
// @Tags examinees
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.PostExamineeInput true "新增常用体检人"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/examinees [post]
func (c *userController) PostExaminee(ctx *gin.Context) {
	var input dto.PostExamineeInput
	if err := util.ParseRequest(ctx, &input); err != nil {
		util.Log.Errorf("参数绑定失败, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}
	userId := ctx.GetInt64("userId")
	id, err := c.service.SaveExaminee(userId, &input)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"user_id": userId}).Errorf("创建常用体检人出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
		return
	}
	middleware.ResponseSuccess(ctx, dto.ResourceID{Id: id})
}

// ListExaminee godoc
// @Summary 个人中心->常用信息
// @Description 获取常用体检人列表
// @Tags examinees
// @Produce json
// @Param  token header string true "用户token"
// @Success 200 {object} middleware.Response{data=[]dto.ListExamineeOutputEle}
// @Router /users/examinees [get]
func (c *userController) ListExaminee(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	output, err := c.service.FindAllExaminees(userId)
	if err != nil {
		util.Log.Errorf("获取用户常用体检人列表出错, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, output)
	}
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
func (c *userController) GETUserProfile(ctx *gin.Context) {
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

// PostUserAddr godoc
// @Summary 新增用户收件地址
// @Description 新增地址
// @Tags addrs
// @Accept  json
// @Produce  json
// @Param token header string true "用户token"
// @Param body body dto.CreateUserAddrInput true "新增用户收件地址"
// @Success 200 {object} middleware.Response{data=dto.ResourceID} "success"
// @Router /users/addrs [post]
func (c *userController) PostUserAddr(ctx *gin.Context) {
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
func (c *userController) PutUserAddr(ctx *gin.Context) {
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
