package main

import (
	"os"

	"mk-api/server/router"
)

func main() {
	server := router.InitRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	_ = server.Run("0.0.0.0:" + port)

}
