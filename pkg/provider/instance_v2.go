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
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"strings"
)

type InstancesV2 struct {
}

func (i *InstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	klog.Info(fmt.Sprintf("InstanceExists providerID:%v", node.Spec.ProviderID))
	nodeInfo, err := client.clientSet.CoreV1().Nodes().Get(ctx, node.Spec.ProviderID, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return true, err
	}
	if nodeInfo != nil && nodeInfo.Labels != nil && nodeInfo.Labels[consts.LabelInstanceType] != "" {
		instanceTypeValue := nodeInfo.Labels[consts.LabelInstanceType]
		list := strings.Split(instanceTypeValue, ".")
		instanceType := instanceTypeValue
		if len(list) > 2 {
			instanceType = list[1]
		}
		switch instanceType {
		case consts.InstanceTypeEcs, consts.InstanceTypeBms, consts.InstanceTypeExternal:
			return true, nil
		}
		klog.Warningf("node %s (providerId:%s) with invalid label: %=%,should been deleted,but exists in kubernetes", node.Name, node.Spec.ProviderID, consts.LabelInstanceType, instanceTypeValue)
		return true, nil
	}
	resp, err := api.NodeCCMInit(consts.ClusterId, node.Spec.ProviderID, "")
	if err != nil {
		return false, err
	}
	if resp.Data.PrivateIp == "" {
		klog.Warningf("can not find providerId for node %s,should been deleted,but exists in kubernetes", node.Name)
		return true, nil
	}
	return true, nil
}

func (i *InstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	return false, errors.New("not implemented")
}

func (i *InstancesV2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	return nil, nil
}
