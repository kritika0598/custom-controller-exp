package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	applister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type controller struct {
	clientset      kubernetes.Interface
	deplister      applister.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, depInformer appsinformer.DeploymentInformer) *controller {
	c := &controller{
		clientset:      clientset,
		deplister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handleDel,
		},
	)

	return c
}

func handleAdd(obj interface{}) {
	fmt.Println("Add was called")
}

func handleDel(obj interface{}) {
	fmt.Println("Delete was called")
}

func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("Statring controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("Waiting for cache to be synced...")
	}

	wait.Until(c.worker, 1*time.Second, ch)
}

func (c *controller) worker() {

}
