package main

import (
	"os"

	"mk-api/server/middleware"
	"mk-api/server/router"
	"mk-api/server/validator"
)

func main() {
	server := router.InitRouter(
		middleware.Secure(),
		middleware.Options(),
		middleware.Logger(),
	)
	port := os.Getenv("PORT")

	// 注册自定义校验器
	validator.Init()

	if port == "" {
		port = "8080"
	}
	_ = server.Run("0.0.0.0:" + port)

}
