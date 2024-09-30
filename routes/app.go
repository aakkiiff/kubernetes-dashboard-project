package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	configinit "kubernetes-api.com/config_init"
	"kubernetes-api.com/models"
)

func stringPtr(s string) *string {
	return &s
}

func createApp(c *gin.Context) {
	var appModel models.App

	namespaceCreationStatus := "namespace created"
	podCreationStatus := "pod created"
	serviceCreationStatus := "service created"
	ingressCreationStatus := "ingress created"

	err := c.ShouldBindJSON(&appModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	NamespaceName := appModel.Name + "-ns"

	ns := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: NamespaceName,
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      appModel.Name,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"podName": appModel.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  appModel.Name + "-container",
					Image: appModel.Image,
				},
			},
		},
	}

	service := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      appModel.Name,
			Namespace: NamespaceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"podName": appModel.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:      appModel.Name + "-ingress",
			Namespace: NamespaceName,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/use-regex":      "true",
				"nginx.ingress.kubernetes.io/rewrite-target": "/$1",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: stringPtr("nginx"),
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/" + appModel.Name,
									PathType: func() *networkingv1.PathType {
										pathType := networkingv1.PathTypePrefix
										return &pathType
									}(),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: appModel.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	clientset := configinit.Initialize_config()

	_, err = configinit.Initialize_config().CoreV1().Namespaces().Create(context.Background(), ns, v1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			// c.JSON(http.StatusInternalServerError, gin.H{"message": "namspace already exists, app will not be created", "error": err.Error()})
			namespaceCreationStatus = "pod already exists, will not be recreated"

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create namespace", "error": err.Error()})
			return
		}

	}
	CheckNamespaceExists(NamespaceName)
	_, err = clientset.CoreV1().Pods(NamespaceName).Create(context.Background(), pod, v1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			podCreationStatus = "pod already exists, will not be recreated"

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create pod", "error": err.Error()})
			return
		}

	}

	_, err = clientset.CoreV1().Services(NamespaceName).Create(context.Background(), service, v1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			serviceCreationStatus = "service already exists, will not be recreated"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Service", "error": err.Error()})
			return
		}

	}

	_, err = clientset.NetworkingV1().Ingresses(NamespaceName).Create(context.Background(), ingress, v1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			ingressCreationStatus = "ingress already exists, will not be recreated"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ingress", "error": err.Error()})
			return
		}

	}
	c.JSON(http.StatusCreated, gin.H{
		"Status": map[string]string{
			"namespace": namespaceCreationStatus,
			"pod":       podCreationStatus,
			"service":   serviceCreationStatus,
			"ingress":   ingressCreationStatus,
		},
		"Application": map[string]string{
			"pod_name":       appModel.Name,
			"namespace":      NamespaceName,
			"containerImage": appModel.Image,
			"service_name":   appModel.Name,
			"ingress":        appModel.Name + "-ingress",
		},
	})
}
