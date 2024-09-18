package provider

import (
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	cloudprovider "k8s.io/cloud-provider"
	"log"
	"time"
)

var _ cloudprovider.Interface = (*Cloud)(nil)

var client *Client

func init() {
	client = initClient()
	cloudprovider.RegisterCloudProvider(consts.ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		return &Cloud{}, nil
	})
}

type Client struct {
	clientSet *kubernetes.Clientset
	informer  cache.SharedIndexInformer
}

func initClient() *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("newCloud:: Failed to create kubernetes config: %v", err))
	}
	clientSet, err := kubernetes.NewForConfig(config)

	nodeInformerFactory := informers.NewSharedInformerFactory(clientSet, time.Second)
	nodeInformer := nodeInformerFactory.Core().V1().Nodes().Informer()
	if err = nodeInformer.AddIndexers(cache.Indexers{
		"spec.providerID": func(obj interface{}) ([]string, error) {
			node, ok := obj.(*v1.Node)
			if !ok {
				return nil, fmt.Errorf("object is not a node")
			}
			return []string{node.Spec.ProviderID}, nil
		},
	}); err != nil {
		panic(fmt.Errorf("can not add indexer %s", err.Error()))
	}
	go nodeInformer.Run(make(chan struct{}))
	go nodeInformerFactory.Start(make(chan struct{}))
	nodeInformerFactory.WaitForCacheSync(make(chan struct{}))
	return &Client{
		clientSet: clientSet,
		informer:  nodeInformer,
	}
}

type Cloud struct {
}

func (cloud *Cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	return
}

func (cloud *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	return &LoadBalancer{
		clientSet: clientSet,
	}, true
}

func (cloud *Cloud) Instances() (cloudprovider.Instances, bool) {
	return &Instances{}, true
}

func (cloud *Cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	return &InstancesV2{clientSet: clientSet}, false
}

func (cloud *Cloud) Zones() (cloudprovider.Zones, bool) {
	return &Zones{}, false
}

func (cloud *Cloud) Clusters() (cloudprovider.Clusters, bool) {
	return &Clusters{}, false
}

func (cloud *Cloud) Routes() (cloudprovider.Routes, bool) {
	return &Routes{}, false
}

func (cloud *Cloud) ProviderName() string {
	return consts.ProviderName
}

func (cloud *Cloud) HasClusterID() bool {
	return true
}

//func (cloud *Cloud) Name() string {
//	return "cdscloud"
//}
