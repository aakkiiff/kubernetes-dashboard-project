package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/node", getNodes)

	server.GET("/namespace", getNamespaces)
	server.POST("/namespace", createNamespace)
	server.DELETE("/namespace", deleteNamespace)

	server.GET("/pods/:ns", getPods)
	server.POST("/pod", createPod)
	server.DELETE("/pod", deletePod)

}
