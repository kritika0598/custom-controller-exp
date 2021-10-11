package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/rest"
	"time"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/Users/kritika/.kube/config", "location to your kube config file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error %s building from flags\n\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error %s  Getting InClusterConfig\n\n", err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s, creating client set \n\n", err.Error())
	}

	ch := make(chan struct{})

	informers := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	c := newController(clientset, informers.Apps().V1().Deployments())

	informers.Start(ch)
	c.run(ch)

	fmt.Println(informers)
}
