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
	)

	// 注册自定义校验器
	validator.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	_ = server.Run("0.0.0.0:" + port)

}
