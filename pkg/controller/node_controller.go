package controller

import (
	"context"
	"errors"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	commoneks "github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/eks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	genericcontrollermanager "k8s.io/controller-manager/app"
	"k8s.io/controller-manager/controller"
	"k8s.io/klog/v2"
	"log"
)

type NodeController struct {
	clientSet *kubernetes.Clientset
}

func NewNodeController() NodeController {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	return NodeController{clientSet: clientSet}
}

func (n *NodeController) StartNodeControllerWrapper(initContext app.ControllerInitContext, completedConfig *cloudcontrollerconfig.CompletedConfig, cloud cloudprovider.Interface) app.InitFunc {

	return func(ctx context.Context, controllerContext genericcontrollermanager.ControllerContext) (controller.Interface, bool, error) {
		return PullNodes(ctx)
	}
}

// PullNodes 设置节点的标签
func PullNodes(ctx context.Context) (controller.Interface, bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Error("list nodes failed")
	}
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
		details, err := eks.DescribeNodeDetails(consts.ClusterId, node.GetName())
		if err != nil {
			return nil, true, nil
		}
		flag, err := UpdateNode(&node, details.Data)
		if err != nil {
			return nil, true, nil
		}
		if flag {
			clientSet.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
		}
	}
	return nil, true, nil
}

// UpdateNode 设置节点的污点
func UpdateNode(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) (bool, error) {
	if detail == nil || node == nil {
		return false, errors.New("invalid node")
	}
	labelFlag := UpdateNodeLabels(node, detail)
	taintFlag := UpdateNodeTaints(node, detail)
	return labelFlag && taintFlag, nil
}

// UpdateNodeLabels 更新节点标签
func UpdateNodeLabels(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) bool {
	if len(detail.Labels) == 0 {
		return false
	}
	labels := make(map[string]string)
	if len(node.Labels) > 0 {
		for key, value := range node.Labels {
			labels[key] = value
		}
	}
	for i := 0; i < len(detail.Labels); i++ {
		label := detail.Labels[i]
		labels[label.Key] = label.Value
	}
	node.Labels = labels
	return true
}

// UpdateNodeTaints 修改节点的污点
func UpdateNodeTaints(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) bool {
	taints := make([]v1.Taint, 0, 0)
	taintMap := make(map[string]v1.Taint)
	if len(detail.Taints) == 0 {
		return false
	}
	for i := 0; i < len(node.Spec.Taints); i++ {
		taint := node.Spec.Taints[i]
		taintMap[taint.Key] = v1.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: taint.Effect,
		}
	}
	for i := 0; i < len(detail.Taints); i++ {
		taint := detail.Taints[i]
		taintMap[taint.Key] = v1.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: v1.TaintEffect(taint.Effect),
		}
	}
	for _, value := range taintMap {
		taints = append(taints, value)
	}
	node.Spec.Taints = taints
	return true
}
