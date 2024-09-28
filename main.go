package main

import (
	"github.com/gin-gonic/gin" // Import corev1 for Namespace
	"kubernetes-api.com/routes"
)

func main() {
	server := gin.Default()

	routes.RegisterRoutes(server)
	server.Run(":8080")
}
