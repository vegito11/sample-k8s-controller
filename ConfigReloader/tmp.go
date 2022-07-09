/* package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vegito11/ConfigReloader/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	QName            = "cm_reloader"
	ANNOTATIONS_KEYS = "meta.reloader.sh/cm-name"
)

func getPods(clientset *kubernetes.Clientset) {

	ctx := context.Background()
	deploys, getErr := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v=%v", ANNOTATIONS_KEYS, "common-config"),
	})

	if getErr != nil {
		panic(fmt.Errorf("failed to get deployment with given configmap mount: %v", getErr))
	}

	for _, dep := range deploys.Items {
		fmt.Println(dep.Spec.Selector)
		labelSelector := ""
		for key, val := range dep.Spec.Selector.MatchLabels {
			labelSelector = key + "=" + val + " "
		}
		pods, _ := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})

		for _, pod := range pods.Items {
			fmt.Printf("\n Deleting pod : %v ", pod.Name)
			// clientset.CoreV1().Pods(pod.Namespace.Delete(ctx, pod.Name, metav1.DDeleteOptions{})
		}
	}

}

func main() {

	kubeconfig := flag.String("kubeconfig", "C:\\Users\\Vegito\\.kube\\config", "location to your kubeconfig file")
	clientset, err := utils.GetKubeConn(kubeconfig)

	if err != nil {
		fmt.Printf("Error %s, creating clientset \n ", err.Error())
	}

	getPods(clientset)

}
*/