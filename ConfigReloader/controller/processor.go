package controller

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

	return ns, name, nil

}

/* Main Worker function which takes item from Queue and process it.  */
func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	defer c.queue.Done(item)

	if shutdown {
		return false
	}

	ns, name, err := ValidateDeploy(c.clientset, item)

	// ===================== 2) Process deleted configmap ------------------
	// If the deployment is not present then it is delete event
	if apierrors.IsNotFound(err) {
		fmt.Printf("\n Handle delete event for configmap %s", name)
	} else if name != "" {
		// ====================== 1) Process updated configmap -------------
		err = c.reloadCM(ns, name)
		if err != nil {
			fmt.Printf("\n Error restarting %v pod %s", name, err.Error())
			return false
		}
	} else if err != nil {
		fmt.Printf("\n Error while processing item %v", err.Error())
		return false
	}

	return true
}
