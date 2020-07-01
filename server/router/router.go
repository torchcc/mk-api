package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mk-api/deployment"
	"mk-api/server/controller"
	"mk-api/server/middleware"

	"mk-api/docs"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// TODO 这里的参数可以考虑zk配置
	docs.SwaggerInfo.Title = "迈康体检网微信服务号 api"

	docs.SwaggerInfo.Description = "迈康体检网微信服务号 API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = getHost()
	docs.SwaggerInfo.BasePath = ""

	router := gin.Default()
	router.Use(middlewares...)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// users
	userRouteGroup := router.Group("/users")
	userRouteGroup.Use(
		middleware.MobileBoundRequired(),
	)

	{
		controller.UserRegister(userRouteGroup)
	}

	// wechat
	weChatRouteGroup := router.Group("/wx")
	{
		controller.WeChatRegister(weChatRouteGroup)
	}

	// login_register
	loginRegisterRouteGroup := router.Group("/login_register")
	loginRegisterRouteGroup.Use(
		middleware.TokenRequired(),
	)

	{
		controller.LoginRegister(loginRegisterRouteGroup)
	}

	// package_register
	pkgRegisterRouteGroup := router.Group("/pkg")
	pkgRegisterRouteGroup.Use(
		middleware.TokenRequired(),
	)

	{
		controller.PackageRegister(pkgRegisterRouteGroup)
	}

	// cart_register
	cartRegisterRouteGroup := router.Group("/cart")
	cartRegisterRouteGroup.Use(
		middleware.MobileBoundRequired(),
	)
	{
		controller.CartRegister(cartRegisterRouteGroup)
	}

	// order_register
	orderRegisterRouteGroup := router.Group("/orders")
	orderRegisterRouteGroup.Use(
		middleware.MobileBoundRequired(),
	)

	{
		controller.OrderRegister(orderRegisterRouteGroup)
	}

	// pay_register, 出于微信回调， 组路由不加mobile required 验证， 需要的话在子路由添加
	payRegisterRouteGroup := router.Group("/pay")

	{
		controller.PayRegister(payRegisterRouteGroup)
	}

	return router
}

func getHost() (host string) {
	switch deployment.BRANCH {
	case "test", "prod":
		host = "106.53.124.190:8081"
	case "local":
		host = "localhost:8081"
	}
	return
}
