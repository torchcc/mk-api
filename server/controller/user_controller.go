package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		userService    service.UserService = service.NewUserService(userModel)
		userController UserController      = NewUserController(userService)
	)
	router.GET("/", userController.FindAll)
	router.GET("/:id", userController.UserDetail)
}

type UserController interface {
	FindAll(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UserDetail(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

var _ *validator.Validate

func NewUserController(service service.UserService) UserController {
	_ = validator.New()
	return &userController{
		service: service,
	}
}

// GetVideos godoc
// @Summary List existing users
// @Description Get all the existing users
// @Tags users,list
// @Accept  json
// @Produce  json
// @Success 200 {array} model.User
// @Failure 401 {object} middleware.Response
// @Router /users [get]
func (c *userController) FindAll(ctx *gin.Context) {
	data, err := c.service.FindAll()
	if err != nil {
		util.Log.Errorf("查找全部用户出错: err: %s\n", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, data)
	}
}

// CreateVideo godoc
// @Summary Create new users
// @Description Create a new user
// @Tags users,create
// @Accept  json
// @Produce  json
// @Param user body model.User true "Create user"
// @Success 200 {object} middleware.Response
// @Failure 500 {object} middleware.Response
// @Router /users [post]
func (c *userController) Create(ctx *gin.Context) {
	var user model.User
	var err error

	if err = ctx.ShouldBindJSON(&user); err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	id, err := c.service.Save(user)
	if err != nil {
		util.Log.Errorf("创建用户出错: err: %s\n", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{ID: id})
	}
}

// UpdateUser godoc
// @Summary Update users
// @Description Update a single user
// @Tags users
// @Accept  json
// @Produce  json
// @Param  id path int true "User ID"
// @Param user body model.User true "Update user"
// @Success 200 {object} middleware.Response
// @Failure 400 {object} middleware.Response
// @Failure 500 {object} middleware.Response
// @Router /users/{id} [put]
func (c *userController) Update(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	user.ID = id

	if err = c.service.Update(user); err != nil {
		util.Log.WithFields(logrus.Fields{
			"user_id": user.ID,
		}).Errorf("创建用户出错: err: %s\n", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
	} else {
		middleware.ResponseSuccess(ctx, dto.ResourceID{ID: user.ID})
	}
}

// UpdateUser godoc
// @Summary Get users
// @Description get a single user's info
// @Tags users
// @Accept json
// @Produce json
// @Param  id path int true "User ID"
// @Success 200 {object} middleware.Response
// @Router /users/{id} [get]
func (c *userController) UserDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		middleware.ResponseError(ctx, ecode.RequestErr, err)
		return
	}

	user, err := c.service.Retrieve(uint32(id))
	if _, ok := errors.Cause(err).(ecode.Codes); ok {
		util.Log.WithFields(logrus.Fields{
			"user_id": id,
		}).Errorf("获取用户信息出错: err: %s\n", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("服务器内部错误"))
	} else {
		middleware.ResponseSuccess(ctx, user)
	}
}
