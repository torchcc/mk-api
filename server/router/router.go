package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mk-api/server/controller"

	"mk-api/docs"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// TODO 这里的参数可以考虑zk配置
	docs.SwaggerInfo.Title = "迈康体检网微信服务号 api"

	docs.SwaggerInfo.Description = "迈康体检网微信服务号 API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""

	router := gin.Default()
	router.Use(middlewares...)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// users
	userRouteGroup := router.Group("/users")
	// usersRouteGroup.Use(
	// 	middleware.TranslationMiddleware(),
	// )
	{
		controller.UserRegister(userRouteGroup)
	}

	// wechat
	weChatRouteGroup := router.Group("/")
	{
		controller.WeChatRegister(weChatRouteGroup)
	}

	return router
}
