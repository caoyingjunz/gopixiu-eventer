package cmd

import (
	"flag"
	"path/filepath"

	"github.com/elastic/go-elasticsearch/v8"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var eventFlag = InitFlags()

func InitFlags() *EventFlags {
	var eventFlags EventFlags

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&eventFlags.KubeConfig, "kubernetes.kubeConfig", filepath.Join(home, ".kube", "config"), "(可选) kubeconfig 文件的绝对路径")
	} else {
		flag.StringVar(&eventFlags.KubeConfig, "kubernetes.kubeConfig", eventFlags.KubeConfig, "kubeconfig 文件的绝对路径")
	}

	flag.StringVar(&eventFlags.Address, "elasticsearch.address", eventFlags.Address, "(可选) elasticsearch address 地址")
	flag.StringVar(&eventFlags.UserName, "elasticsearch.username", eventFlags.UserName, "(可选) elasticsearch 用户名")
	flag.StringVar(&eventFlags.Password, "elasticsearch.password", eventFlags.UserName, "(可选) elasticsearch 用户名")

	flag.Parse()
	return &eventFlags
}

func InitClient() (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", eventFlag.KubeConfig); err != nil {
			panic(err.Error())
		}
	}
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeclient, nil
}

func InitElasticSearch() (*elasticsearch.Client, error) {
	var err error
	var es *elasticsearch.Client
	cfg := elasticsearch.Config{
		Addresses: []string{
			eventFlag.Address,
		},
		Username: eventFlag.UserName,
		Password: eventFlag.Password,
		//		CACert: cert,
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	_, err = es.Info()
	if err != nil {
		panic(err.Error())
	}
	return es, nil
}
