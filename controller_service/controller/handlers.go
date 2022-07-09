package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *controller) processNewDeploy(ns, name string) error {

	deploy, getErr := c.depLister.Deployments(ns).Get(name)

	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}

	if ns == "kube-sysem" {
		return nil
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
			Annotations: map[string]string{
				ANNOTATIONS[0]["key"]: ANNOTATIONS[0]["value"],
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
			Selector: deploy.Spec.Template.Labels,
		},
	}

	_, err := c.clientset.CoreV1().Services(ns).Create(context.Background(), &svc, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf(" Creating Service %s \n", err.Error())
		return err
	}
	fmt.Printf("\n Created service 🎯 in %v with name %v", deploy.Namespace, deploy.Name)

	return nil
}

func (c *controller) processRemovedDeploy(ns, name string) error {

	svc, getErr := c.clientset.CoreV1().Services(ns).Get(context.Background(), name, metav1.GetOptions{})

	if apierrors.IsNotFound(getErr) {
		fmt.Printf("\n Service with name %s does not exists \n", getErr.Error())
		return getErr
	} else if getErr != nil {
		fmt.Printf("\n Error While Getting the Service %s \n", getErr.Error())
	}

	if value := svc.Annotations[ANNOTATIONS[0]["key"]]; value != ANNOTATIONS[0]["value"] {
		fmt.Printf("\n Skipping Deletion of the Service %s \n", svc.Name)
		return nil
	}

	delErr := c.clientset.CoreV1().Services(ns).Delete(context.Background(), name, metav1.DeleteOptions{})

	if delErr != nil {
		fmt.Printf("\n Sucessfully deleted the Service %s 🪢 \n", svc.Name)
		return nil
	}

	return delErr
}
