package cmd

import (
	"flag"
	"github.com/elastic/go-elasticsearch/v8"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func InitClient() (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选) kubeconfig 文件的绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的绝对路径")
	}
	flag.Parse()

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
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
			"http://10.50.7.126:9200",
		},
		Username: "elastic",
		Password: "elastic",
		//		CACert: cert,
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return es, nil
}