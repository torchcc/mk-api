package main

import (
	"os"

	"mk-api/server/middleware"
	"mk-api/server/router"
)

func main() {
	server := router.InitRouter(middleware.Secure())
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	_ = server.Run("0.0.0.0:" + port)

}
