package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/node", getNodes)
	server.GET("/namespace", getNamespaces)
	server.POST("/namespace", createNamespace)
	server.DELETE("/namespace", deleteNamespace)

	server.POST("pod", createPod)

}
