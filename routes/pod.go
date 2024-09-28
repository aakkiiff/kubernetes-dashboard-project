package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configinit "kubernetes-api.com/config_init"
	"kubernetes-api.com/models"
)

func getPods(c *gin.Context) {
	var podNames []string
	ns := c.Param("ns")
	// Get the Kubernetes clientset
	clientset := configinit.Initialize_config()

	// Create the Pod in the specified namespace
	pods, err := clientset.CoreV1().Pods(ns).List(context.Background(), v1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not get the pods", "error": err.Error()})
		return
	}
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	c.JSON(200, gin.H{"namespace": ns, "pods": podNames})
}

func createPod(c *gin.Context) {
	// Bind the JSON body to the Pod struct
	var podModel models.Pod
	err := c.ShouldBindJSON(&podModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	// Create a Kubernetes Pod object
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      podModel.Name,
			Namespace: podModel.NamespaceName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  podModel.Name + "-container",
					Image: podModel.Image,
				},
			},
		},
	}

	// Get the Kubernetes clientset
	clientset := configinit.Initialize_config()

	// Create the Pod in the specified namespace
	createdPod, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.Background(), pod, v1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create pod", "error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pod created successfully",
		"pod": map[string]string{
			"pod_name":       createdPod.Name,
			"namespace":      createdPod.Namespace,
			"container":      createdPod.Spec.Containers[0].Name,
			"containerImage": createdPod.Spec.Containers[0].Image,
		},
	})
}

func deletePod(c *gin.Context) {
	// Bind the JSON body to the Pod struct
	var podModel models.Pod
	err := c.ShouldBindJSON(&podModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	// Get the Kubernetes clientset
	clientset := configinit.Initialize_config()
	err = clientset.CoreV1().Pods(podModel.NamespaceName).Delete(context.Background(), podModel.Name, v1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not dleete the pod", "error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "pod deleted successfully"})

}
