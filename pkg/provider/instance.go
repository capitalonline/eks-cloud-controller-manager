package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

type Instances struct {
}

func (i *Instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.Info(fmt.Sprintf("NodeAddresses name:%v", name))
	address, err := api.NodeAddresses(consts.ClusterId, "", string(name))
	if err != nil {
		klog.Errorf("查询节点ip失败,err:%s", err.Error())
		return nil, nil
	}
	if len(address) == 0 {
		node, err := client.clientSet.CoreV1().Nodes().Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return node.Status.Addresses, nil
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
	if len(address) == 0 {
		node, err := i.getNodeByByProviderID(providerID)
		if err != nil {
			return nodeAddress, err
		}
		nodeAddress = node.Status.Addresses
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
	node, err := client.clientSet.CoreV1().Nodes().Get(ctx, string(nodeName), metav1.GetOptions{})
	if err != nil {
		klog.Errorf("查询节点%s失败,err:%v", nodeName, err)
		return "", err
	}
	if node != nil && node.Spec.ProviderID != "" {
		return node.Spec.ProviderID, nil
	}
	klog.Warningf("could not find instanceId for node %s ,should been deleted,but exists in kubernetes", nodeName)
	return "", fmt.Errorf("could not find instanceId for node %s", nodeName)
}

func (i *Instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.Info(fmt.Sprintf("InstanceType name:%v", name))

	resp, err := api.NodeCCMInit(consts.ClusterId, "", string(name))
	if err != nil {
		klog.Errorf("通过openapi查询节点%s失败,err:%v", string(name), err)
		return "", err
	}
	for j := 0; j < len(resp.Data.Labels); j++ {
		label := resp.Data.Labels[j]
		if label.Key == consts.LabelInstanceType {
			return label.Value, nil
		}
	}
	node, err := client.clientSet.CoreV1().Nodes().Get(ctx, string(name), metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}
	if node != nil && node.Labels != nil && node.Labels[consts.LabelInstanceType] != "" {
		instanceType := node.Labels[consts.LabelInstanceType]
		return instanceType, nil
	}
	return "external", nil
}

func (i *Instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.Info(fmt.Sprintf("InstanceTypeByProviderID providerID:%v", providerID))
	var node *v1.Node
	resp, err := api.NodeCCMInit(consts.ClusterId, providerID, "")
	if err != nil {
		klog.Errorf("通过openapi查询节点%s失败,err:%v", providerID, err)
		return "", err
	}
	for j := 0; j < len(resp.Data.Labels); j++ {
		label := resp.Data.Labels[j]
		if label.Key == consts.LabelInstanceType {
			return label.Value, nil
		}
	}
	node, err = i.getNodeByByProviderID(providerID)
	if err != nil {
		return "", err
	}
	if node != nil && node.Labels[consts.LabelInstanceType] != "" {
		instanceType := node.Labels[consts.LabelInstanceType]
		return instanceType, nil
	}
	return "external", nil
}

func (i *Instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return errors.New("not implemented")
}

func (i *Instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *Instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.Info(fmt.Sprintf("InstanceExistsByProviderID providerID:%v", providerID))
	// 如果providerId包含
	if len(providerID) == 0 {
		return true, errors.New("providerID为空")
	}
	node, err := i.getNodeByByProviderID(providerID)
	if err != nil {
		return false, err
	}
	if node.Name != "" && node.Labels != nil && node.Labels[consts.LabelInstanceType] != "" {
		instanceTypeValue := node.Labels[consts.LabelInstanceType]
		list := strings.Split(instanceTypeValue, ".")
		instanceType := instanceTypeValue
		if len(list) > 2 {
			instanceType = list[1]
		}
		switch instanceType {
		case consts.InstanceTypeEcs, consts.InstanceTypeBms, consts.InstanceTypeExternal:
			return true, nil
		}
		klog.Warningf("node %s (providerId:%s) with invalid label: %=%,should been deleted,but exists in kubernetes", node.Name, providerID, consts.LabelInstanceType, instanceTypeValue)
		return true, nil
	}
	address, err := api.NodeCCMInit(consts.ClusterId, providerID, "")
	if err != nil {
		klog.Errorf("通过openapi查询节点%s失败,err:%v", providerID, err)
		return true, err
	}
	switch address.Data.Status {
	// 需要删除
	case consts.NodeStatusDeleted:
		//klog.Warningf("node %v is deleted by server", providerID)
		klog.Warningf("node %q (providerId:%s) deleted from eks-server,but exists in kubernetes", node.Name, providerID)
		return true, nil
	default:
	}
	return true, nil
}

func (i *Instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, errors.New("not implemented")
}

func (i *Instances) getNodeByByProviderID(providerID string) (*v1.Node, error) {
	time.Sleep(time.Second * 2)
	var node *v1.Node
	nodes, err := client.informer.GetIndexer().ByIndex("spec.providerID", providerID)
	if err != nil && !apierrors.IsNotFound(err) {
		return node, err
	}
	if len(nodes) < 1 {
		return node, errors.New("can't find node by providerID")
	}
	node, ok := nodes[0].(*v1.Node)
	if !ok {
		return node, fmt.Errorf("nodes[0] is not a node, %v", node)
	}
	return node, nil
}
