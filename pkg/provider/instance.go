package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

type Instances struct {
}

func (i *Instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.Info(fmt.Sprintf("NodeAddresses name:%v", name))
	address, err := api.NodeAddresses(consts.ClusterId, "", string(name))
	if err != nil {
		// TODO 加日志
		klog.Errorf("")
		return nil, nil
	}
	nodeAddress := make([]v1.NodeAddress, 0, len(address))
	for _, item := range address {
		if item == "" {
			continue
		}
		nodeAddress = append(nodeAddress, v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: item,
		})
	}
	return nodeAddress, nil
}

func (i *Instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	klog.Info(fmt.Sprintf("NodeAddressesByProviderID    providerID:%v", providerID))
	if providerID == "" {
		return nil, errors.New("providerID can not be empty")
	}
	address, err := api.NodeAddresses(consts.ClusterId, providerID, "")
	if err != nil {
		// TODO 加日志
		return nil, nil
	}
	nodeAddress := make([]v1.NodeAddress, 0, len(address))
	for _, item := range address {
		if item == "" {
			continue
		}
		nodeAddress = append(nodeAddress, v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: item,
		})
	}
	return nodeAddress, nil
}

func (i *Instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.Info(fmt.Sprintf("InstanceID nodeName:%v", nodeName))
	resp, err := api.NodeCCMInit(consts.ClusterId, "", string(nodeName))
	if err != nil {
		return "", err
	}
	// TODO node_id为空处理
	return resp.Data.NodeId, nil
}

func (i *Instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.Info(fmt.Sprintf("InstanceID name:%v", name))

	resp, err := api.NodeCCMInit(consts.ClusterId, "", string(name))
	if err != nil {
		return "", err
	}
	for j := 0; j < len(resp.Data.Labels); j++ {
		label := resp.Data.Labels[j]
		if label.Key == "node.kubernetes.io/instance.type" {
			return label.Value, err
		}
	}
	return "", nil
}

func (i *Instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.Info(fmt.Sprintf("InstanceTypeByProviderID providerID:%v", providerID))
	resp, err := api.NodeCCMInit(consts.ClusterId, providerID, "")
	if err != nil {
		return "", err
	}
	for j := 0; j < len(resp.Data.Labels); j++ {
		label := resp.Data.Labels[j]
		if label.Key == "node.kubernetes.io/instance.type" {
			return label.Value, err
		}
	}
	return "ecs", nil
}

func (i *Instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return errors.New("not implemented")
}

func (i *Instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *Instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.Info(fmt.Sprintf("InstanceExistsByProviderID providerID:%v", providerID))
	address, err := api.NodeCCMInit(consts.ClusterId, providerID, "")
	if err != nil {
		return true, err
	}
	// TODO NodeId
	if address.Data.NodeId == "" {
		return false, nil
	}
	return true, nil
}

func (i *Instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, errors.New("not implemented")
}
