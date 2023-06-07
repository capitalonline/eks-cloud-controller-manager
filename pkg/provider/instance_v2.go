package provider

import (
	"context"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

type InstancesV2 struct {
}

func (i *InstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	return false, nil
}

func (i *InstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	return false, nil
}

func (i *InstancesV2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	return nil, nil
}
