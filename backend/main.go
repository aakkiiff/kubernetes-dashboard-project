package main

import (
	"github.com/gin-gonic/gin"
	"kubernetes-api.com/routes"
)

func main() {
	server := gin.Default()

	routes.RegisterRoutes(server)
	server.Run(":8080")
}
