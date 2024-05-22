package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"strings"
)

type InstancesV2 struct {
}

func (i *InstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	klog.Info(fmt.Sprintf("InstanceExists providerID:%v", node.Spec.ProviderID))
	if strings.Contains(node.Spec.ProviderID, "external-node") {
		return true, nil
	}
	resp, err := api.NodeCCMInit(consts.ClusterId, node.Spec.ProviderID, "")
	if err != nil {
		return false, nil
	}
	if resp.Data.PrivateIp == "" {
		return false, nil
	}
	return true, nil
}

func (i *InstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	return false, errors.New("not implemented")
}

func (i *InstancesV2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	return nil, nil
}
