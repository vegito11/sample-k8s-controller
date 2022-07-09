package controller

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	coreinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	QName            = "cm_reloader"
	ANNOTATIONS_KEYS = "meta.reloader.sh/cm-name"
)

type controller struct {
	clientset     kubernetes.Interface
	cmLister      corelisters.ConfigMapLister
	cmCacheSynced cache.InformerSynced
	queue         workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface, cmInformer coreinformer.ConfigMapInformer) *controller {

	c := &controller{
		clientset:     clientset,
		cmLister:      cmInformer.Lister(),
		cmCacheSynced: cmInformer.Informer().HasSynced,
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), QName),
	}

	cmInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: c.handleDel,
			UpdateFunc: c.handleUpdate,
		},
	)

	return c
}

func (c *controller) enqueCM(obj interface{}) (string, error) {

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		fmt.Printf("\n Error getting key from cache %s \n", err.Error())
		return "", err
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		fmt.Printf("\n Error while spliting key into namespace and name")
		return "", err
	}

	for _, nsm := range []string{"kube-system", "ingress-nginx",
		"kube-public", "kube-node-lease"} {
		if ns == nsm {
			return "", errors.New("skip: restricted namespace cm")
		}
	}

	c.queue.Add(obj)
	return name, nil

}

func (c *controller) handleUpdate(new, obj interface{}) {

	if nm, err := c.enqueCM(obj); err == nil {
		fmt.Printf("\n Update Method Called for configmap %v ➕ ", nm)
	}
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("\n Delete Method Called ❌")
	// c.queue.Add(obj)
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("\n Starting Config Reloader Controller ⚡")

	// Informer maintain local cache , so we have to wait till cache to be synced
	// at least for one time
	if !cache.WaitForCacheSync(ch, c.cmCacheSynced) {
		fmt.Print("\n Waiting for cache to be synced \n")
	}
	// It will call c.worker function after every 1 second
	go wait.Until(c.worker, 1*time.Second, ch)
	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}
