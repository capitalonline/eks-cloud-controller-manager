package provider

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"io"
	"k8s.io/client-go/informers"
	cloudprovider "k8s.io/cloud-provider"
)

func init() {
	cloudprovider.RegisterCloudProvider(consts.ProviderName, func(io.Reader) (cloudprovider.Interface, error) {
		return &Cloud{}, nil
	})
}

type Cloud struct {
}

func (cloud *Cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	clientSet := clientBuilder.ClientOrDie("cds-cloud-provider")
	sharedInformer := informers.NewSharedInformerFactory(clientSet, 0)
	sharedInformer.Start(stop)
	sharedInformer.WaitForCacheSync(stop)
	return
}

func (cloud *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return &LoadBalancer{}, true
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
	return &Routes{}, true
}

func (cloud *Cloud) ProviderName() string {
	return consts.ProviderName
}

func (cloud *Cloud) HasClusterID() bool {
	return false
}
