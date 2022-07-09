package controller

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const QName string = "service_expose"

var ANNOTATIONS = []map[string]string{
	{"key": "managed.by", "value": "expose-controller"},
}

type controller struct {
	clientset      kubernetes.Interface
	depLister      appslisters.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface, depInformer appsinformer.DeploymentInformer) *controller {

	c := &controller{
		clientset:      clientset,
		depLister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), QName),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)

	return c
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("Add Method Called ➕")
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("Delete Method Called ❌")
	c.queue.Add(obj)
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println(" Starting Controller ⚡")

	// Informer maintain local cache , so we have to wait till cache to be synced
	// at least for one time
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("Waiting for cache to be synced \n")
	}
	// It will call c.worker function after every 1 second
	go wait.Until(c.worker, 1*time.Second, ch)
	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}
