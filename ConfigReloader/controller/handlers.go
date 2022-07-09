package controller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *controller) reloadCM(ns, cmname string) error {

	ctx := context.Background()
	deploys, getErr := c.clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v=%v", ANNOTATIONS_KEYS, cmname),
	})

	if getErr != nil {
		fmt.Printf(" \n failed to get deployment with given configmap mount: %v", getErr)
		return getErr
	}

	for ind, dep := range deploys.Items {

		fmt.Printf("\n %v) Processing Deployment : %v ", ind+1, dep.Name)
		labelSelector := ""
		for key, val := range dep.Spec.Selector.MatchLabels {
			labelSelector = key + "=" + val + " "
		}
		pods, _ := c.clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})

		for _, pod := range pods.Items {

			fmt.Printf("\n Deleting pod : %v ", pod.Name)
			if delErr := c.clientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); delErr != nil {
				fmt.Printf("\n %v) Error while deleting the pod ", delErr.Error())
				return delErr
			}
		}
	}

	return nil
}
