package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)


type Controller struct {
	indexer cache.Indexer
	queue workqueue.RateLimitingInterface
	informer cache.Controller
	es *elasticsearch.Client
}

func NewController(indexer cache.Indexer, queue workqueue.RateLimitingInterface, informer cache.Controller, es *elasticsearch.Client) *Controller {
	return &Controller{
		queue: queue,
		indexer: indexer,
		informer: informer,
		es: es,
	}
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}


func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer  c.queue.Done(key)

	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		fmt.Printf("Event %s does not exists anymore\n", key)
	} else {
		eventResp := &EventResp{
			NAMESPACE: obj.(*v1.Event).InvolvedObject.Namespace,
			LastTimestamp: obj.(*v1.Event).LastTimestamp.UTC(),
			Type: obj.(*v1.Event).Type,
			Reason: obj.(*v1.Event).Reason,
			Object: fmt.Sprintf("%s/%s",obj.(*v1.Event).InvolvedObject.Kind,obj.(*v1.Event).InvolvedObject.Name),
			Message: obj.(*v1.Event).Message,
		}
		data, err := json.Marshal(eventResp)
		if err != nil {
			log.Fatalf("Error marshaling document: %s", err)
			return err
		}
		req := esapi.IndexRequest{
			Index:      "event",
			DocumentID: string(obj.(*v1.Event).UID),
			Body:       bytes.NewReader(data),
			Refresh:    "true",
		}
		res, err := req.Do(context.Background(), c.es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
			return err
		}
		defer res.Body.Close()
		fmt.Printf("插入数据完成  %s\n", obj.(*v1.Event).InvolvedObject.Name)
	}
	return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	defer c.queue.ShutDown()
	klog.Infof("Starting Event Controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh

	klog.Info("Stopping Pod controller")
}