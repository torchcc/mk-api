package main

import (
	"os"

	"mk-api/server/router"
)

func main() {
	server := router.InitRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}
	server.Run(":" + port)

}
