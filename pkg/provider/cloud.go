package provider

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cloudprovider "k8s.io/cloud-provider"
	"log"
)

var _ cloudprovider.Interface = (*Cloud)(nil)

func init() {
	cloudprovider.RegisterCloudProvider(consts.ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		return &Cloud{}, nil
	})
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
	return &InstancesV2{}, false
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
