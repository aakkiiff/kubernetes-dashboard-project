package routes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configinit "kubernetes-api.com/config_init"
)

func getNodes(c *gin.Context) {
	var nodeNames []string
	nodes, err := configinit.Initialize_config().CoreV1().Nodes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println("can not list node names", err.Error())
		return
	}
	for _, nodes := range nodes.Items {
		// fmt.Println(nodes.Name)
		nodeNames = append(nodeNames, nodes.Name)
	}
	c.JSON(200, gin.H{"nodes": nodeNames})

}
