package provider

import (
	"context"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/eks"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Instances struct {
}

func (i *Instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	address, err := eks.NodeAddresses(consts.ClusterId, "", string(name))
	if err != nil {
		return nil, err
	}
	nodeAddress := make([]v1.NodeAddress, 0, len(address))
	for _, item := range address {
		nodeAddress = append(nodeAddress, v1.NodeAddress{
			Address: item,
		})
	}
	return nodeAddress, nil
}

func (i *Instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	address, err := eks.NodeAddresses(consts.ClusterId, providerID, "")
	if err != nil {
		return nil, err
	}
	nodeAddress := make([]v1.NodeAddress, 0, len(address))
	for _, item := range address {
		nodeAddress = append(nodeAddress, v1.NodeAddress{
			Address: item,
		})
	}
	return nil, nil
}

func (i *Instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {

	return "", nil
}

func (i *Instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	return "", nil
}

func (i *Instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "", nil
}

func (i *Instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return nil
}

func (i *Instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return "", nil
}

func (i *Instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, nil
}

func (i *Instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, nil
}
