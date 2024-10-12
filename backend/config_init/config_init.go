package configinit

import (
	"flag"
	"log"
	"os"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
	once      sync.Once
)

func Initialize_config() *kubernetes.Clientset {
	once.Do(func() {
		config, err := rest.InClusterConfig()
		if err != nil {
			kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "absolute path to the kubeconfig file")
			flag.Parse()

			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				log.Fatalf("Error building kubeconfig: %v", err)
			}
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %v", err)
		}
	})

	return clientset
}
