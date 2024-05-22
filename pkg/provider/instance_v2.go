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
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"time"
)

type InstancesV2 struct {
	clientSet *kubernetes.Clientset
}

func (i *InstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	klog.Info(fmt.Sprintf("InstanceExists providerID:%v", node.Spec.ProviderID))
	_, err := i.clientSet.CoreV1().Nodes().Get(ctx, node.Spec.ProviderID, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return true, err
	}
	if !apierrors.IsNotFound(err) {
		if _, ok := node.Labels[consts.LabelExternalNode]; ok {
			return true, nil
		}
		if time.Now().UTC().Sub(node.CreationTimestamp.Time.UTC()) < time.Minute*15 {
			return true, nil
		}
		return false, nil
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
