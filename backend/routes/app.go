package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
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
func int32Ptr(i int32) *int32 { return &i }

func createApp(c *gin.Context) {
	var appModel models.App

	namespaceCreationStatus := "namespace created"
	deploymentCreationStatus := "deployment created"
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

	// pod := &corev1.Pod{
	// 	ObjectMeta: v1.ObjectMeta{
	// 		Name:      appModel.Name,
	// 		Namespace: NamespaceName,
	// 		Labels: map[string]string{
	// 			"podName": appModel.Name,
	// 		},
	// 	},
	// 	Spec: corev1.PodSpec{
	// 		Containers: []corev1.Container{
	// 			{
	// 				Name:  appModel.Name + "-container",
	// 				Image: appModel.Image,
	// 				Ports: []corev1.ContainerPort{
	// 					{
	// 						ContainerPort: int32(appModel.ContainerPort),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	deployment := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      appModel.Name,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"app": appModel.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appModel.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": appModel.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  appModel.Name + "-container",
							Image: appModel.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          appModel.Name,
									ContainerPort: int32(appModel.ContainerPort),
								},
							},
						},
					},
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
				"app": appModel.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(appModel.ContainerPort),
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
									Path: "/" + appModel.Name + "/?(.*)",
									PathType: func() *networkingv1.PathType {
										pathType := networkingv1.PathTypeImplementationSpecific
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
			namespaceCreationStatus = "namespace already exists, will not be recreated"

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create namespace", "error": err.Error()})
			return
		}

	}
	CheckNamespaceReady(NamespaceName)

	// _, err = clientset.CoreV1().Pods(NamespaceName).Create(context.Background(), pod, v1.CreateOptions{})
	// if err != nil {
	// 	if errors.IsAlreadyExists(err) {
	// 		podCreationStatus = "pod already exists, will not be recreated"

	// 	} else {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create pod", "error": err.Error()})
	// 		return
	// 	}

	// }

	_, err = clientset.AppsV1().Deployments(NamespaceName).Create(context.Background(), deployment, v1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			deploymentCreationStatus = "Deployment already exists, will not be recreated"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create deployment", "error": err.Error()})
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
			"namespace":  namespaceCreationStatus,
			"deployment": deploymentCreationStatus,
			"service":    serviceCreationStatus,
			"ingress":    ingressCreationStatus,
		},
		"Application": map[string]string{
			"deployment_name": appModel.Name,
			"namespace":       NamespaceName,
			"containerImage":  appModel.Image,
			"service_name":    appModel.Name,
			"ingress":         appModel.Name + "-ingress",
		},
	})
}

func updateApp(c *gin.Context) {
	var appModel models.App
	status := "no updates!"
	err := c.ShouldBindJSON(&appModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	NamespaceName := appModel.Name + "-ns"
	clientset := configinit.Initialize_config()

	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), NamespaceName, v1.GetOptions{})

	if err != nil {
		c.JSON(http.StatusBadRequest, "application do not exist! Please create an application first")
		return
	}

	// pod, err := clientset.CoreV1().Pods(NamespaceName).Get(context.Background(), appModel.Name, v1.GetOptions{})
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, "pod do not exist! Someone might have deleted it! Please re-apply the application first")
	// 	return
	// }

	deployment, err := clientset.AppsV1().Deployments(NamespaceName).Get(context.Background(), appModel.Name, v1.GetOptions{})

	if err != nil {
		c.JSON(http.StatusBadRequest, "deployment do not exist! Someone might have deleted it! Please re-apply the application first")
		return
	}

	service, err := clientset.CoreV1().Services(NamespaceName).Get(context.Background(), appModel.Name, v1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, "service do not exist! Someone might have deleted it! Please re-apply the application first")
		return
	}

	_, err = clientset.NetworkingV1().Ingresses(NamespaceName).Get(context.Background(), appModel.Name+"-ingress", v1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, "ingress do not exist! Someone might have deleted it! Please re-apply the application first")
		return
	}

	if appModel.Image != deployment.Spec.Template.Spec.Containers[0].Image {
		deployment.Spec.Template.Spec.Containers[0].Image = appModel.Image

		_, err = clientset.AppsV1().Deployments(NamespaceName).Update(context.Background(), deployment, v1.UpdateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update the deployment", "error": err.Error()})
			return
		}
		status = "application image updated!"

	}

	if appModel.ContainerPort != int(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort) {

		deployment, err := clientset.AppsV1().Deployments(NamespaceName).Get(context.Background(), appModel.Name, v1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusBadRequest, "deployment do not exist! Someone might have deleted it! Please re-apply the application first")
			return
		}

		deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(appModel.ContainerPort)
		service.Spec.Ports[0].TargetPort = intstr.FromInt(appModel.ContainerPort)

		_, err = clientset.AppsV1().Deployments(NamespaceName).Update(context.Background(), deployment, v1.UpdateOptions{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update the deployment", "error": err.Error()})
			return
		}

		_, err = clientset.CoreV1().Services(NamespaceName).Update(context.Background(), service, v1.UpdateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update the service", "error": err.Error()})
			return
		}

		if status == "application image updated!" {
			status = "application image & port updated!"
		} else {
			status = "application port updated"
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": status})
}
