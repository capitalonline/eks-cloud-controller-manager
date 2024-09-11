package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	commoneks "github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	"github.com/go-ping/ping"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"reflect"
	"strings"
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
	//metricList, err := n.metricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	//if err != nil {
	//	klog.Info("获取指标列表失败，err:", err)
	//	return err
	//}
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
	//for _, metric := range metricList.Items {
	//	node := nodeSet[metric.Name]
	//	load, err := n.CalculateLoad(metric, node)
	//	if err != nil {
	//		continue
	//	}
	//	request.NodeList = append(request.NodeList, commoneks.ModifyClusterLoadReqNode{
	//		NodeId:   node.Spec.ProviderID,
	//		NodeName: node.Name,
	//		//Cpu: &commoneks.ResourceInfo{
	//		//	Usage:    load.Cpu.Usage,
	//		//	Limits:   load.Cpu.Limits,
	//		//	Requests: load.Cpu.Requests,
	//		//},
	//		//Memory: &commoneks.ResourceInfo{
	//		//	Usage:    load.Mem.Usage,
	//		//	Limits:   load.Mem.Limits,
	//		//	Requests: load.Mem.Requests,
	//		//},
	//		Status: load.Status,
	//	})
	//}
	for _, node := range nodeList.Items {
		//node := nodeSet[metric.Name]
		//load, err := n.CalculateLoad(metric, node)
		//if err != nil {
		//	continue
		//}
		status := consts.NodeStatusNotReady
		if NodeReady(node) {
			status = consts.NodeStatusReady
		}

		request.NodeList = append(request.NodeList, commoneks.ModifyClusterLoadReqNode{
			NodeId:   node.Spec.ProviderID,
			NodeName: node.Name,
			//Cpu: &commoneks.ResourceInfo{
			//	Usage:    load.Cpu.Usage,
			//	Limits:   load.Cpu.Limits,
			//	Requests: load.Cpu.Requests,
			//},
			//Memory: &commoneks.ResourceInfo{
			//	Usage:    load.Mem.Usage,
			//	Limits:   load.Mem.Limits,
			//	Requests: load.Mem.Requests,
			//},
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

func (n *NodeController) CalculateLoad(metric v1beta1.NodeMetrics, node v1.Node) (commoneks.NodeLoad, error) {
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
		//status     = "NotReady"
		status = consts.NodeStatusNotReady
	)
	if NodeReady(node) {
		status = consts.NodeStatusReady
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
	return commoneks.NodeLoad{
		Cpu: commoneks.ResourceInfo{
			Usage:    int64(cpuUsage * 100),
			Limits:   int64(float64(limitCpu) / float64(allCpu) * 100),
			Requests: int64(float64(requestCpu) / float64(allCpu) * 100),
		},
		Mem: commoneks.ResourceInfo{
			Usage:    int64(memoryUsage * 100),
			Limits:   int64(float64(limitMem) / float64(allMem) * 100),
			Requests: int64(float64(requestMem) / float64(allMem) * 100),
		},
		Status: status,
	}, nil
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

func (n *NodeController) ListenNodes(ctx context.Context) {

	factory := informers.NewSharedInformerFactory(n.clientSet, time.Second)
	informer := factory.Core().V1().Events().Informer()
	stopCh := make(chan struct{})
	factory.Core().V1().Events().Lister()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event, _ := obj.(*v1.Event)
			switch event.Reason {
			case consts.EventNodeNotReady:
				n.NotifyNodeDown(ctx, event)
			case consts.EventNodeReady:
				n.NotifyNodeReady(ctx, event)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			event, _ := newObj.(*v1.Event)
			switch event.Reason {
			case consts.EventNodeNotReady:
				n.NotifyNodeDown(ctx, event)
			case consts.EventNodeReady:
				n.NotifyNodeReady(ctx, event)
			}
		},
		//DeleteFunc: func(obj interface{}) {
		//	event, _ := obj.(*v1.Event)
		//	klog.Errorf("delete event %s", event.Reason)
		//},
	},
	)

	factory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		klog.Errorf("同步事件失败")
		return
	}
	fmt.Println("informer run")
	go informer.Run(stopCh)
	<-ctx.Done()
}

func (n *NodeController) NotifyNodeReady(ctx context.Context, event *v1.Event) {
	if event.InvolvedObject.Kind != consts.ResourceKindNode || event.Namespace != v1.NamespaceDefault {
		return
	}
	node, err := n.clientSet.CoreV1().Nodes().Get(ctx, event.InvolvedObject.Name, metav1.GetOptions{})
	if err != nil || node == nil || !NodeReady(*node) {
		return
	}

	var request = commoneks.NewModifyClusterLoadRequest()
	request.ClusterId = consts.ClusterId
	request.NodeList = []commoneks.ModifyClusterLoadReqNode{
		{
			NodeId:   node.Spec.ProviderID,
			Status:   consts.NodeStatusReady,
			NodeName: node.Name,
		},
	}
	_, err = api.ModifyClusterLoad(request)
	if err != nil {
		klog.Errorf("notify NodeReady failed，err:%v", err)
		return
	}
	recordName := fmt.Sprintf("ccm-%s-down", node.Name)
	if err = n.clientSet.CoreV1().Events(v1.NamespaceDefault).Delete(ctx, recordName, metav1.DeleteOptions{}); err != nil && !kerrors.IsNotFound(err) {
		klog.Errorf("delete event %s err:%s", recordName, err.Error())
	}
	return
}

func (n *NodeController) NotifyNodeDown(ctx context.Context, event *v1.Event) {
	if event.Source.Component != "node-controller" || event.Namespace != v1.NamespaceDefault {
		//klog.Errorf("event source is not node-controller")
		return
	}
	node, err := n.clientSet.CoreV1().Nodes().Get(ctx, event.InvolvedObject.Name, metav1.GetOptions{})
	if _, ok := node.Labels[consts.NodeRoleMaster]; !ok {
		return
	}
	if err != nil || node == nil || NodeHealth(*node) {
		data, _ := json.Marshal(node)
		klog.Infof("err: %v,node ready %v, Unreachable:%v,nodehealth:%v, node:%s, ", err, NodeReady(*node), NodeUnreachable(*node), NodeHealth(*node), string(data))
		return
	}
	klog.Info("node is unhealth")
	// worker节点忽略

	recordName := fmt.Sprintf("ccm-%s-down", node.Name)
	record, err := n.clientSet.CoreV1().Events(v1.NamespaceDefault).Get(ctx, recordName, metav1.GetOptions{})
	if err != nil && !kerrors.IsNotFound(err) {
		klog.Errorf("get ccm event record failed, err:%s", err.Error())
		return
	}
	klog.Infof("record %s not found", recordName)
	if record != nil && record.Name == recordName {
		data, _ := json.Marshal(record)
		klog.Infof("record %s is not nil,record: %s", recordName, string(data))
		return
	}
	ip := node.Name
	if len(strings.Split(node.Name, "-")) > 1 {
		ip = strings.Split(node.Name, "-")[1]
	}

	var req = commoneks.NewModifyClusterLoadRequest()
	req.ClusterId = consts.ClusterId
	req.NodeList = []commoneks.ModifyClusterLoadReqNode{
		{
			NodeId:   node.Spec.ProviderID,
			Status:   consts.NodeStatusNotReady,
			NodeName: node.Name,
		},
	}
	if _, err = api.ModifyClusterLoad(req); err != nil {
		klog.Errorf("report node %s status failed,err:%v", node.Name, err)
		return
	}
	request := commoneks.NewSendAlarmRequest()
	request.Theme = consts.K8sMetricAlarmTheme
	request.NodeId = node.Spec.ProviderID
	request.ClusterId = consts.ClusterId
	request.Metric = consts.AlarmMetricMasterDown
	request.Source = consts.AlarmSource
	request.AlarmMsg = fmt.Sprintf("集群%s节点NotReady", node.Name)
	request.Tags = []interface{}{}
	request.Keyword = node.Name
	request.Value = ip
	resp, err := api.NotifyMasterDown(request)
	if err != nil || resp == nil || resp.Code != consts.EksRequestSuccess {
		klog.Error(fmt.Sprintf("send alarm error：%v, resp:%v", err, resp))
		return
	}
	record = &v1.Event{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      recordName,
			Namespace: v1.NamespaceDefault,
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      event.InvolvedObject.Kind,
			Namespace: event.InvolvedObject.Namespace,
			Name:      event.InvolvedObject.Name,
		},
		Reason:  consts.EventNodeNotReady,
		Message: fmt.Sprintf("master %s down", node.Name),
		Source: v1.EventSource{
			Component: "cloud-ctroller-manager",
			Host:      ip,
		},
		FirstTimestamp: event.FirstTimestamp,
		LastTimestamp:  metav1.Time{Time: time.Now()},
		Count:          1,
		Type:           v1.EventTypeNormal,
		EventTime:      event.EventTime,
	}
	_, err = n.clientSet.CoreV1().Events(v1.NamespaceDefault).Create(ctx, record, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		klog.Errorf("crate event %s err:%s", recordName, err.Error())
	}
	n.clientSet.CoreV1().ConfigMaps(consts.NameSpaceKubeSystem).Get(ctx, "kubeadm-config", metav1.GetOptions{})
	klog.Info("notify down success")
}

func NodeReady(node v1.Node) bool {
	if len(node.Status.Conditions) < 1 {
		return false
	}
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func NodeUnreachable(node v1.Node) bool {
	if len(node.Spec.Taints) < 1 {
		return false
	}
	for _, taint := range node.Spec.Taints {
		if taint.Key == v1.TaintNodeUnreachable {
			return true
		}
	}
	return false
}

func NodeHealth(node v1.Node) bool {
	ip, err := NodeIP(node)
	if err != nil {
		return false
	}
	if !NodeReady(node) && NodeUnreachable(node) {
		return false
	}
	pinger := ping.New(ip)
	pinger.Count = 10
	pinger.Timeout = time.Second * 2
	pinger.SetPrivileged(true)
	if err = pinger.Run(); err != nil {
		return false
	}
	if pinger.Statistics().PacketLoss > 50 {
		return false
	}
	return true
}

func NodeIP(node v1.Node) (string, error) {
	if len(node.Status.Addresses) < 1 {
		return "", errors.New("node's status.addresses is nil")
	}
	for _, addr := range node.Status.Addresses {
		if addr.Type == consts.NodeAddrTypeInternalIP {
			return addr.Address, nil
		}
	}
	return "", errors.New("can not find ip")
}

func (n *NodeController) Update(ctx context.Context) error {
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
		// TODO 批量查
		details, err := api.NodeCCMInit(consts.ClusterId, node.Spec.ProviderID, "")
		if err != nil {
			return err
		}
		// id为空可能是节点添加后任务流还未回调ccm就查询了，需要跳过
		if details == nil || details.Data.NodeId == "" {
			continue
		}
		switch details.Data.Status {
		case consts.NodeStatusDeleted:
			//需要ccm主动触发删除该节点
			//if err := n.clientSet.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{}); err != nil {
			//	klog.Errorf("unable to delete node %q: %v", node.Name, err)
			//}
			klog.Warningf("node %q (providerId:%s) deleted from eks-server,but exists in kubernetes", node.Name, node.Spec.ProviderID)
		case consts.NodeStatusRunning:
			oldNode := node.DeepCopy()
			if err := UpdateNode(&node, details.Data); err != nil {
				return fmt.Errorf("update node %s failed with error %v", node.Name, err)
			}

			//flag, err := UpdateNode(&node, details.Data)
			//if err != nil {
			//	return err
			//}
			if !reflect.DeepEqual(node.Labels, oldNode.Labels) ||
				!reflect.DeepEqual(node.Spec.Taints, oldNode.Spec.Taints) ||
				!reflect.DeepEqual(node.Annotations, oldNode.Annotations) {
				_, err = n.clientSet.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
				if err != nil {
					klog.Errorf("更新节点失败，err:%v", err)
				}
			}

			//if flag {
			//	_, err = n.clientSet.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
			//	if err != nil {
			//		klog.Errorf("更新节点失败，err:%v", err)
			//	}
			//}
		default:
		}
	}
	return nil
}

// UpdateNode 设置节点的污点
func UpdateNode(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) error {
	if detail == nil || node == nil {
		return errors.New("invalid node")
	}
	//labelFlag := UpdateNodeLabels(node, detail)
	//taintFlag := UpdateNodeTaints(node, detail)
	//annotationFlag := UpdateNodeAnnotations(node, detail)
	UpdateNodeLabels(node, detail)
	UpdateNodeTaints(node, detail)
	UpdateNodeAnnotations(node, detail)
	return nil
}

// UpdateNodeLabels 更新节点标签
func UpdateNodeLabels(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) {
	if len(detail.Labels) == 0 {
		return
	}
	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}
	for i := 0; i < len(detail.Labels); i++ {
		label := detail.Labels[i]
		node.Labels[label.Key] = label.Value
	}
}

// UpdateNodeTaints 修改节点的污点
func UpdateNodeTaints(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) {

	taints := make([]v1.Taint, 0, 0)
	taintMap := make(map[string]v1.Taint)
	if len(detail.Taints) == 0 {
		return
	}
	if len(node.Spec.Taints) == 0 {
		node.Spec.Taints = taints
		return
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
}

// UpdateNodeAnnotations 修改节点的污点
func UpdateNodeAnnotations(node *v1.Node, detail *commoneks.NodeCCMInitResponseData) {
	annotations := make(map[string]string)
	if len(detail.Annotations) == 0 {
		return
	}
	if len(node.Annotations) == 0 {
		node.Annotations = annotations
		return
	}
	for i := 0; i < len(detail.Annotations); i++ {
		annotation := detail.Annotations[i]
		node.Annotations[annotation.Key] = annotation.Value
	}
}
