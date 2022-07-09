package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/vegito11/ConfigReloader/controller"
	"github.com/vegito11/ConfigReloader/utils"

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

	c := controller.NewController(clientset, informers.Core().V1().ConfigMaps())
	informers.Start(ch)
	c.Run(ch)
	fmt.Println(informers)
}
