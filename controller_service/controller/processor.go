package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func ValidateDeploy(clientset kubernetes.Interface, item interface{}) (string, string, error) {

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("\n Error getting key from cache %s \n", err.Error())
		return "", "", err
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		fmt.Printf("\n Error while spliting key into namespace and name")
		return "", "", err
	}

	deploy, err := clientset.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		fmt.Printf("\n  Deployment with %s is not found in ns %s deployment ", name, ns)
		return ns, name, err
	}

	if value := deploy.Annotations[ANNOTATIONS[0]["key"]]; value == ANNOTATIONS[0]["value"] {
		return ns, name, nil
	}

	fmt.Printf("\n Deployment %v Missing the required Annotation %v : %v",
		name, ANNOTATIONS[0]["key"], ANNOTATIONS[0]["value"])
	return "", "", nil

}

/* Main Worker function which takes item from Queue and process it.  */
func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	defer c.queue.Done(item)

	if shutdown {
		return false
	}

	ns, name, err := ValidateDeploy(c.clientset, item)

	// ===================== 2) Process deleted deployment ------------------
	// If the deployment is not present then it is delete event
	if apierrors.IsNotFound(err) {
		fmt.Printf("\n  Handle delete event for deployment %s", name)
		c.processRemovedDeploy(ns, name)
	} else if name != "" {
		// ====================== 1) Process newly Added deployment -------------
		err = c.processNewDeploy(ns, name)
		if err != nil {
			fmt.Printf("\n  Error syncing %v deployment %s", name, err.Error())
			return false
		}
	} else if err != nil {
		fmt.Printf("\n  Error while processing item %v", err.Error())
		return false
	}

	return true
}
