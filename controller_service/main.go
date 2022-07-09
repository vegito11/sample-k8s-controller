package main

import (
	"controller_service/controller"
	"controller_service/utils"
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/informers"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "C:\\Users\\Vegito\\.kube\\config", "location to your kubeconfig file")
	clientset, err := utils.GetKubeConn(kubeconfig)

	if err != nil {
		fmt.Printf("Error %s, creating clientset \n ", err.Error())
	}

	ch := make(chan struct{})
	informers := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	if err != nil {
		fmt.Printf("Getting informer factory %s \n", err.Error())
	}

	c := controller.NewController(clientset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.Run(ch)
	fmt.Println(informers)
}
