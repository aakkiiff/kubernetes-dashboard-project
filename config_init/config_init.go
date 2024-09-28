package configinit

import (
	"flag"
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
	once      sync.Once
)

// Intt initializes the Kubernetes client only once
func Initialize_config() *kubernetes.Clientset {
	once.Do(func() {
		// Define the flag for kubeconfig
		kubeconfig := flag.String("kubeconfig", "/home/akif/.kube/config", "location to your kubeconfig file")

		// Parse the command-line flags
		flag.Parse()

		// Build the Kubernetes config using the provided kubeconfig path
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %v", err)
		}

		// Create a Kubernetes client from the config
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %v", err)
		}
	})

	return clientset
}
