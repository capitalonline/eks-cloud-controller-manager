package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	commoneks "github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"time"
)

type NodeController struct {
	clientSet     *kubernetes.Clientset
	metricsClient *metrics.Clientset
}

func (n *NodeController) Validate() error {
	return nil
}

func NewNodeController() NodeController {
	klog.Info("NewNodeController")
	//config, err := clientcmd.RESTConfigFromKubeConfig([]byte(s))
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return NodeController{clientSet: clientSet, metricsClient: metricsClient}
}

// CollectPlayLoad 获取集群节点的负载信息
func (n *NodeController) CollectPlayLoad(ctx context.Context) error {
	metricList, err := n.metricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Info("获取指标列表失败，err:", err)
		return err
	}
	// 查询所有节点信息，获取余量信息
	nodeList, err := n.clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Info("获取节点列表失败，err:", err)
		return err
	}
	nodeSet := make(map[string]v1.Node)
	for _, node := range nodeList.Items {
		nodeSet[node.Name] = node
	}
	var request = commoneks.NewModifyClusterLoadRequest()
	request.ClusterId = consts.ClusterId
	request.NodeList = make([]commoneks.ModifyClusterLoadReqNode, 0)
	for _, metric := range metricList.Items {

		node := nodeSet[metric.Name]
		usage := metric.Usage
		cpuUsage := float64(usage.Cpu().MilliValue()) / float64(node.Status.Allocatable.Cpu().MilliValue())
		memoryUsage := float64(usage.Memory().MilliValue()) / float64(node.Status.Allocatable.Memory().MilliValue())
		//cpuRequests := float64(usage.Cpu().MilliValue())/float64(node.Status)
		var (
			requestCpu int64
			requestMem int64
			limitCpu   int64
			limitMem   int64
			allCpu     = node.Status.Allocatable.Cpu().MilliValue()
			allMem     = node.Status.Allocatable.Memory().MilliValue()
			status     = "NotReady"
		)

		for i := 0; i < len(node.Status.Conditions); i++ {
			condition := node.Status.Conditions[i]
			if condition.Type == "Ready" {
				if condition.Status == "True" {
					status = "Ready"
				}
			}
		}
		pods, _ := n.clientSet.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name)})
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				requestCpu += container.Resources.Requests.Cpu().MilliValue()
				requestMem += container.Resources.Requests.Memory().MilliValue()
				limitCpu += container.Resources.Limits.Cpu().MilliValue()
				limitMem += container.Resources.Limits.Cpu().MilliValue()
			}
		}
		request.NodeList = append(request.NodeList, commoneks.ModifyClusterLoadReqNode{
			NodeId:   node.Spec.ProviderID,
			NodeName: node.Name,
			Cpu: commoneks.ResourceInfo{
				Usage:    int64(cpuUsage * 100),
				Limits:   int64(float64(limitCpu) / float64(allCpu) * 100),
				Requests: int64(float64(requestCpu) / float64(allCpu) * 100),
			},
			Memory: commoneks.ResourceInfo{
				Usage:    int64(memoryUsage * 100),
				Limits:   int64(float64(limitMem) / float64(allMem) * 100),
				Requests: int64(float64(requestMem) / float64(allMem) * 100),
			},
			Status: status,
		})
	}
	_, err = api.ModifyClusterLoad(request)
	if err != nil {
		klog.Info("同步节点负载失败，err:", err)
		return err
	}
	klog.Info("更新节点负载成功")
	return nil
}

func (n *NodeController) Run(ctx context.Context) error {
	klog.Info("开始运行run")
	ticker := time.NewTicker(time.Minute * 3)
	for {
		select {
		case <-ticker.C:
			err := n.Update(ctx)
			if err != nil {
				klog.Infoln(err)
			}
			err = n.CollectPlayLoad(ctx)
			if err != nil {
				klog.Infoln(err)
			}
		case <-ctx.Done():
			klog.Info("程序退出")
			return nil
		}
	}
}

func (n *NodeController) Update(ctx context.Context) error {
	klog.Info("开始获取节点信息")
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	//}
	//clientSet, err := kubernetes.NewForConfig(config)
	nodes, err := n.clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Error("list nodes failed")
	}
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
		details, err := api.NodeCCMInit(consts.ClusterId, node.Spec.ProviderID, "")
		if err != nil {
			return err
		}
		flag, err := UpdateNode(&node, details.Data)
		if err != nil {
			return err
		}
		if flag {
			_, err = n.clientSet.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
			fmt.Println(err)
		}
		// 需要ccm主动触发删除该节点
		if details != nil && details.Data.NodeId == "" {
			if err := n.clientSet.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{}); err != nil {
				klog.Errorf("unable to delete node %q: %v", node.Name, err)
			}
		}
	}
	return nil
}

// UpdateNode 设置节点的污点
func UpdateNode(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) (bool, error) {
	if detail == nil || node == nil {
		return false, errors.New("invalid node")
	}
	labelFlag := UpdateNodeLabels(node, detail)
	taintFlag := UpdateNodeTaints(node, detail)
	annotationFlag := UpdateNodeAnnotations(node, detail)
	klog.Info(fmt.Sprintf("更新节点%s,labelFlag:%v,taintFlag:%v,annotationFlag:%v", node.Name, labelFlag, taintFlag, annotationFlag))
	return labelFlag || taintFlag || annotationFlag, nil
}

// UpdateNodeLabels 更新节点标签
func UpdateNodeLabels(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) bool {
	klog.Info(fmt.Sprintf("更新节点%s的标签,原标签：%v, eks标签：%v", node.Name, node.Labels, detail.Labels))
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
func UpdateNodeTaints(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) bool {
	klog.Info(fmt.Sprintf("更新节点%s的污点,原污点：%v, eks污点：%v", node.Name, node.Spec.Taints, detail.Taints))

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

// UpdateNodeAnnotations 修改节点的污点
func UpdateNodeAnnotations(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) bool {
	klog.Info(fmt.Sprintf("更新节点%s的注释,原注释：%v, eks注释：%v", node.Name, node.Spec.Taints, detail.Taints))
	annotations := make(map[string]string)
	if len(detail.Annotations) == 0 {
		return false
	}
	for i := 0; i < len(detail.Annotations); i++ {
		annotation := detail.Annotations[i]
		annotations[annotation.Key] = annotation.Value
	}
	if len(node.Annotations) == 0 {
		node.Annotations = annotations
		return true
	}
	for k, v := range node.Annotations {
		annotations[k] = v
	}
	node.Annotations = annotations
	return true
}
