package main

import (
	"github.com/darianJmy/event-collect/cmd"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)



func main() {
	client, err := cmd.InitClient()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	es, err := cmd.InitElasticSearch()
	if err != nil {
		klog.Fatalf("Error building elasticsearch client: %s", err.Error())
	}
	// 创建 event ListWatcher
	eventListWatcher := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "events", v1.NamespaceAll, fields.Everything())
	// 创建队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	// 创建 indexer、informer
	indexer, informer := cache.NewIndexerInformer(eventListWatcher, &v1.Event{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	}, cache.Indexers{})

	// 实例化 controller
	controller := cmd.NewController(indexer, queue, informer, es)

	// 启动 controller
	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(1, stopCh)

	select {}
}

