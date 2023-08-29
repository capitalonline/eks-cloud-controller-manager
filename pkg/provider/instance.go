package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type Instances struct {
	clientSet *kubernetes.Clientset
}

func (i *Instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.Info(fmt.Sprintf("NodeAddresses name:%v", name))
	address, err := api.NodeAddresses(consts.ClusterId, "", string(name))
	if err != nil {
		klog.Errorf("查询节点ip失败,err:%s", err.Error())
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
		klog.Errorf("查询节点ip失败,err:%s", err.Error())
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
		klog.Errorf("通过openapi查询节点%s失败,err:%v", string(nodeName), err)
		return "", err
	}
	//  NodeId
	if resp.Data.NodeId != "" {
		return resp.Data.NodeId, nil
	}
	// NodeId 为空可能是因为节点的状态改变，直接返回providerId
	node, err := i.clientSet.CoreV1().Nodes().Get(ctx, string(nodeName), metav1.GetOptions{})
	if err != nil {
		klog.Errorf("查询节点%s失败,err:%v", nodeName, err)
		return "", err
	}
	if node != nil {
		return node.Spec.ProviderID, nil
	}
	return "", nil
}

func (i *Instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.Info(fmt.Sprintf("InstanceID name:%v", name))

	resp, err := api.NodeCCMInit(consts.ClusterId, "", string(name))
	if err != nil {
		klog.Errorf("通过openapi查询节点%s失败,err:%v", string(name), err)
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

func (i *Instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.Info(fmt.Sprintf("InstanceTypeByProviderID providerID:%v", providerID))
	resp, err := api.NodeCCMInit(consts.ClusterId, providerID, "")
	if err != nil {
		klog.Errorf("通过openapi查询节点%s失败,err:%v", providerID, err)
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
		klog.Errorf("通过openapi查询节点%s失败,err:%v", providerID, err)
		return true, err
	}
	switch address.Data.Status {
	// 这三种状态是需要删除的
	case consts.NodeStatusFailed, consts.NodeStatusError, consts.NodeStatusDeleted:
		return false, nil
	default:
	}
	return true, nil
}

func (i *Instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, errors.New("not implemented")
}
