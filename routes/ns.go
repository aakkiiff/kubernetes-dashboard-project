package routes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configinit "kubernetes-api.com/config_init"
	"kubernetes-api.com/models"
)

func getNamespaces(c *gin.Context) {
	var namespaceNames []string
	namespaces, err := configinit.Initialize_config().CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println("can not list namespace names", err.Error())
		return
	}
	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}
	c.JSON(200, gin.H{"namespaces": namespaceNames})
}

func createNamespace(c *gin.Context) {
	var namespaceName models.Namespace
	err := c.ShouldBindJSON(&namespaceName)
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	// config_init := configinit.Initialize_config()

	nsName := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespaceName.Name,
		},
	}
	_, err = configinit.Initialize_config().CoreV1().Namespaces().Create(context.Background(), nsName, v1.CreateOptions{})
	if err != nil {
		c.JSON(400, gin.H{"message": "could not create ns", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "namespace " + namespaceName.Name + " created successfully"})
}

func deleteNamespace(c *gin.Context) {
	var namespaceName models.Namespace
	err := c.ShouldBindJSON(&namespaceName)
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	// fmt.Println(namespaceName.Name)
	err = configinit.Initialize_config().CoreV1().Namespaces().Delete(context.Background(), namespaceName.Name, v1.DeleteOptions{})
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid request, could not delete the namespace", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "namespace " + namespaceName.Name + " deleted successfully"})

}
