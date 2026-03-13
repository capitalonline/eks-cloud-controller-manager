package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/lb"
	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

const EKS = "eks"

const (
	AnnotationLbId  = "service.beta.kubernetes.io/cds-select-slb-id"
	AnnotationLbEip = "service.beta.kubernetes.io/cds-select-slb-eip-addr"
	AnnotationLbVip = "service.beta.kubernetes.io/cds-select-slb-vip-addr"

	AnnotationLbNetwork   = "service.beta.kubernetes.io/cds-load-balancer-network"
	AnnotationLbProtocol  = "service.beta.kubernetes.io/cds-load-balancer-protocol"
	AnnotationLbType      = "service.beta.kubernetes.io/cds-load-balancer-types"
	AnnotationLbSpec      = "service.beta.kubernetes.io/cds-load-balancer-specification"
	AnnotationLbBandwidth = "service.beta.kubernetes.io/cds-load-balancer-bandwidth"

	AnnotationLbAlgorithm           = "service.beta.kubernetes.io/cds-load-balancer-algorithm"
	AnnotationLbSubjectId           = "service.beta.kubernetes.io/cds-load-balancer-subject-id"
	AnnotationLbBillingMethod       = "service.beta.kubernetes.io/cds-load-balancer-billingmethod"
	AnnotationLbBillingMethodConfId = "service.beta.kubernetes.io/cds-load-balancer-billingmethod-conf-id"

	LabelNodeRegionCode = "topology.kubernetes.io/region"
	LabelNodeAzCode     = "node.kubernetes.io/node.az-code"
)

const (
	LbNetTypeWan    = "wan"
	LbNetTypeWanLan = "wan_lan"

	LbBillingMethodCostPay = "0" // 按需计费
	BandwidthShared        = "shared"

	DefaultBillingType    = "number"
	FlowDemandBillingType = "flow_demand"

	LbTaskSuccess = "success"
	LbTakError    = "error"

	LBSpecStandard = "standard" // 标准型
	LBSpecHigh     = "high"     // 高阶型
	LBSpecSuper    = "super"    // 超强型
	LBSpecExtreme  = "extreme"  // 至强型

	LBSpecNameStandard = "标准型"
	LBSpecNameHigh     = "高阶型"
	LBSpecNameSuper    = "超强型"
	LBSpecNameExtreme  = "至强型"

	IpTypeInternal = "InternalIP"

	UpdateListenFull  = "full"  // 全量更新
	UpdateListenExact = "exact" // 精确更新

	AllNetwork     = "all"
	PublicNetwork  = "public"
	PrivateNetwork = "private"
	EIP            = "eip"
	LanVip         = "lan_vip"

	LbAlgorithmRr   = "rr"
	LbAlgorithmWrr  = "wrr"
	LbAlgorithmHash = "conhash"

	SlbBuilding = "创建中"
	SlbRunning  = "正常"
)

var lbSpecMap = map[string]string{
	LBSpecStandard: LBSpecNameStandard,
	LBSpecHigh:     LBSpecNameHigh,
	LBSpecSuper:    LBSpecNameSuper,
	LBSpecExtreme:  LBSpecNameExtreme,
}

var lbConfMap = map[string]string{
	LBSpecStandard: "slb.v1.mini",
	LBSpecHigh:     "slb.v1.small",
	LBSpecSuper:    "slb.v1.medium",
	LBSpecExtreme:  "slb.v1.large",
}

var SLBNotFound error = errors.New("slb not found")

// 服务参数结构体
type serviceParams struct {
	subjectId           int
	lbType              int64
	lbSpec              string
	billingMethod       string
	billingMethodConfId int
	lbBandwidth         int64
	lbEip               string
	lbVip               string
	selectSLB           string
	lbNetworkType       string
	ingressStatusIpMap  map[string]bool
}

type LoadBalancer struct {
	clientSet *kubernetes.Clientset
}

// GetLoadBalancer 查询lb
func (l *LoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	response, err := l.describeLbInstance(ctx, service)
	if err != nil || response == nil {
		// k8s在删除节点之后会查一遍slb，确认是否被删除
		if errors.Is(err, SLBNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	// 接口请求返回异常
	slb := response.Data
	ingresses := make([]v1.LoadBalancerIngress, 0, len(slb.VipList))
	for i := 0; i < len(slb.VipList); i++ {
		vipInfo := slb.VipList[i]
		ingresses = append(ingresses, v1.LoadBalancerIngress{
			IP: vipInfo.Vip,
		})
	}
	return &v1.LoadBalancerStatus{
		Ingress: ingresses,
	}, true, err
}

// GetLoadBalancerName 获取lb名称
func (l *LoadBalancer) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	response, err := l.describeLbInstance(ctx, service)
	if err != nil || response == nil {
		return ""
	}
	return response.Data.SlbName
}

// EnsureLoadBalancer 创建lb
func (l *LoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	// 验证服务配置
	if err := l.validateService(service); err != nil {
		return nil, err
	}

	// 获取或创建SLB实例
	slbInfo, err := l.getOrCreateSlb(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create SLB: %w", err)
	}

	if slbInfo.SlbStatus == SlbBuilding {
		return nil, errors.New("slb is building, please wait")
	}

	if slbInfo.SlbStatus != SlbRunning {
		return nil, fmt.Errorf("slb status: %s", slbInfo.SlbStatus)
	}

	if service.Spec.ExternalTrafficPolicy == v1.ServiceExternalTrafficPolicyTypeLocal {
		nodes, err = l.makeLocalLbListen(ctx, service)
		if err != nil {
			return nil, fmt.Errorf("failed to make local lb listen: %w", err)
		}
	}

	// 更新负载均衡监听器
	vipList, err := l.updateLbListen(ctx, service, nodes, slbInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update load balancer: %w", err)
	}

	// 获取最终SLB状态
	return l.makeLoadBalancerStatus(vipList)
}

func (l *LoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	resp, err := l.describeLbInstance(ctx, service)
	if err != nil || resp == nil {
		return fmt.Errorf("UpdateLoadBalancer failed, describe SLB error: %w", err)
	}

	if resp.Data.SlbStatus == SlbBuilding {
		return errors.New("[UpdateLoadBalancer] slb is building, please wait")
	}

	if resp.Data.SlbStatus != SlbRunning {
		return fmt.Errorf("[UpdateLoadBalancer] slb status: %s", resp.Data.SlbStatus)
	}

	if service.Spec.ExternalTrafficPolicy == v1.ServiceExternalTrafficPolicyTypeLocal {
		nodes, err = l.makeLocalLbListen(ctx, service)
		if err != nil {
			return fmt.Errorf("failed to make local lb listen: %w", err)
		}
	}
	_, err = l.updateLbListen(ctx, service, nodes, &resp.Data)

	return err
}

func (l *LoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	response, err := l.describeLbInstance(ctx, service)
	if err != nil || response == nil {
		if errors.Is(err, SLBNotFound) {
			return nil
		}
		return fmt.Errorf("EnsureLoadBalancerDeleted failed, describe SLB error: %w", err)
	}

	return l.clearLbListen(ctx, service, &response.Data)
}

// validateService 验证服务配置的有效性
func (l *LoadBalancer) validateService(service *v1.Service) error {
	if service.Spec.SessionAffinity != v1.ServiceAffinityNone {
		return errors.New("SessionAffinity is not supported currently, only support 'None' type")
	}
	return nil
}

// getOrCreateSlb 获取现有SLB或创建新的SLB
func (l *LoadBalancer) getOrCreateSlb(ctx context.Context, service *v1.Service) (*lb.DescribeVpcSlbResponseSlbInfo, error) {
	// 查询现有SLB实例
	describeResp, err := l.describeLbInstance(ctx, service)
	if err != nil && !errors.Is(err, SLBNotFound) {
		return nil, err
	}

	// 如果SLB不存在则创建新实例
	if describeResp == nil || len(describeResp.Data.SlbId) < 1 {
		return l.createSlb(ctx, service)
	}

	// 返回已存在的SLB ID
	return &describeResp.Data, nil
}

// getLoadBalancerStatus 获取负载均衡器的状态信息
func (l *LoadBalancer) makeLoadBalancerStatus(ips []*lb.DescribeVpcSlbResponseVipInfo) (*v1.LoadBalancerStatus, error) {
	// 构建负载均衡器入口状态
	ingresses := make([]v1.LoadBalancerIngress, 0, len(ips))
	for _, vipInfo := range ips {
		ingresses = append(ingresses, v1.LoadBalancerIngress{
			IP: vipInfo.Vip,
		})
	}

	return &v1.LoadBalancerStatus{
		Ingress: ingresses,
	}, nil
}

func (l *LoadBalancer) createSlb(ctx context.Context, service *v1.Service) (*lb.DescribeVpcSlbResponseSlbInfo, error) {
	// 输入验证
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	// 解析和验证服务参数
	params, err := l.parseServiceParams(service)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service parameters: %w", err)
	}

	if params.selectSLB != "" {
		return nil, SLBNotFound
	}

	// 获取可用区信息
	regionCode, azCode, err := l.getAvailableZone(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available zone: %w", err)
	}
	taskId := ""
	switch params.lbNetworkType {
	case PrivateNetwork:
		taskId, err = l.createStandardSlb(regionCode, azCode, service, params)
		if err != nil {
			return nil, err
		}

	default:
		taskId, err = l.createPackageSlb(azCode, service, params)
		if err != nil {
			return nil, err
		}
	}

	// 等待任务完成
	if err = l.describeTask(taskId); err != nil {
		return nil, err
	}

	describeResp, err := l.describeLbInstance(ctx, service)
	if err != nil || describeResp == nil {
		return nil, fmt.Errorf("failed to describe SLB: %w", err)
	}

	return &describeResp.Data, nil
}

func (l *LoadBalancer) createPackageSlb(azCode string, service *v1.Service, params *serviceParams) (string, error) {
	// 获取计费方案
	billingSchemeId, err := l.getBillingSchemeId(azCode, params.lbSpec)
	if err != nil {
		return "", fmt.Errorf("failed to get billing scheme: %w", err)
	}

	// 获取带宽计费方案
	bandwidthBillingSchemeId, err := l.getBandwidthBillingSchemeId(azCode, params.billingMethod, params.billingMethodConfId)
	if err != nil {
		return "", fmt.Errorf("failed to get bandwidth billing scheme: %w", err)
	}
	if bandwidthBillingSchemeId == "" {
		return "", fmt.Errorf("billing method '%s' not found", params.billingMethod)
	}

	// 构建创建请求
	request := l.buildCreateSlbRequest(service, params, azCode, billingSchemeId, bandwidthBillingSchemeId)

	// 发起创建请求
	response, err := api.PackageCreateSlb(request)
	if err != nil || response == nil {
		return "", fmt.Errorf("API call PackageCreateSlb failed: %w", err)
	}

	if response.Code != consts.LbRequestSuccess {
		return "", fmt.Errorf("create package lb failed, code: %s, message: %s", response.Code, response.Message)
	}

	return response.TaskId, nil
}

func (l *LoadBalancer) createStandardSlb(regionCode, azCode string, service *v1.Service, params *serviceParams) (string, error) {
	req := lb.NewStandardCreateSlbRequest()
	req.RegionCode = regionCode
	req.AvailableZoneCode = azCode
	req.Name = SlbName(service.Name, service.Namespace, string(service.UID))
	req.VpcId = consts.VpcID
	req.NetType = LbNetTypeWanLan

	if params.lbSpec == "" {
		return "", fmt.Errorf("missing required annotation: %s", AnnotationLbSpec)
	}

	confType := lbConfMap[params.lbSpec]
	if confType == "" {
		return "", fmt.Errorf("not fount lb spec conf type '%s'", params.lbSpec)
	}
	req.ConfType = confType
	response, err := api.StandardCreateSlb(req)
	if err != nil || response == nil {
		return "", fmt.Errorf("API call StandardCreateSlb failed: %w", err)
	}

	if response.Code != consts.LbRequestSuccess {
		return "", fmt.Errorf("create standard lb failed, code: %s, message: %s", response.Code, response.Message)
	}
	return response.TaskId, nil
}

func (l *LoadBalancer) makeLocalLbListen(ctx context.Context, service *v1.Service) ([]*v1.Node, error) {
	// 1. 查询 service 对应的 Endpoints
	endpoints, err := l.clientSet.CoreV1().Endpoints(service.Namespace).Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		if k8serr.IsNotFound(err) {
			klog.Infof("endpoints not found for service %s/%s, skip local lb listen", service.Namespace, service.Name)
			return []*v1.Node{}, nil
		}
		return nil, fmt.Errorf("failed to get endpoints for service %s/%s: %w", service.Namespace, service.Name, err)
	}

	// 2. 提取 Endpoints 中的 subsets 信息
	var nodeNames = make(map[string]bool)
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			if address.NodeName != nil {
				nodeNames[*address.NodeName] = true
			}
		}
	}

	// 3. 根据节点名称查询对应的 v1.Node 对象
	var nodes []*v1.Node
	for nodeName, _ := range nodeNames {
		node, e := l.clientSet.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
		if e != nil {
			klog.Warningf("make local Lb listen, failed to get ep node %s: %v", nodeName, err)
			continue
		}
		nodes = append(nodes, node)
	}

	// 4. 返回所有节点信息
	return nodes, nil
}

// 解析服务参数
func (l *LoadBalancer) parseServiceParams(service *v1.Service) (*serviceParams, error) {
	var (
		subjectId int
		m         = make(map[string]bool)
	)

	if len(service.Annotations) == 0 {
		return nil, errors.New("service annotations is null")
	}

	if subjectIdStr := service.Annotations[AnnotationLbSubjectId]; subjectIdStr != "" {
		id, err := strconv.ParseInt(subjectIdStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid subject ID '%s': %w", subjectIdStr, err)
		}
		subjectId = int(id)
	}

	lbTypeStr := service.Annotations[AnnotationLbType]
	if lbTypeStr == "" {
		lbTypeStr = "4"
	}
	lbType, err := strconv.ParseInt(lbTypeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid lb type '%s': %w", lbTypeStr, err)
	}

	lbSpec := service.Annotations[AnnotationLbSpec]

	// 带宽计费方式
	billingMethod := service.Annotations[AnnotationLbBillingMethod]
	switch billingMethod {
	case DefaultBillingType, FlowDemandBillingType:
	default:
		billingMethod = DefaultBillingType
	}

	confId := 0
	billingMethodConfId := service.Annotations[AnnotationLbBillingMethodConfId]
	if billingMethodConfId != "" {
		confId64, _ := strconv.ParseInt(billingMethodConfId, 10, 64)
		confId = int(confId64)
	}

	// 固定带宽
	lbBandwidthStr := service.Annotations[AnnotationLbBandwidth]
	if lbBandwidthStr == "" {
		lbBandwidthStr = "50"
	}
	lbBandwidth, err := strconv.ParseInt(lbBandwidthStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bandwidth '%s': %w", lbBandwidthStr, err)
	}

	networkType := service.Annotations[AnnotationLbNetwork]
	if networkType == "" {
		networkType = PublicNetwork
	}

	lbEipStr := service.Annotations[AnnotationLbEip]
	lbVipStr := service.Annotations[AnnotationLbVip]
	selectSlb := service.Annotations[AnnotationLbId]

	if len(service.Status.LoadBalancer.Ingress) > 0 {
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				m[ingress.IP] = true
			}
		}
		klog.Infof("get lb ingress ip status: %v", m)
	}

	// 验证参数范围
	if lbType <= 0 || lbBandwidth <= 0 {
		return nil, errors.New("lb annotation error")
	}

	return &serviceParams{
		subjectId:           subjectId,
		lbType:              lbType,
		lbSpec:              lbSpec,
		billingMethod:       billingMethod,
		billingMethodConfId: confId,
		lbBandwidth:         lbBandwidth,
		lbEip:               lbEipStr,
		lbVip:               lbVipStr,
		selectSLB:           selectSlb,
		lbNetworkType:       networkType,
		ingressStatusIpMap:  m,
	}, nil
}

func (l *LoadBalancer) getProtocol(service *v1.Service) string {
	if service.Annotations == nil {
		return "TCP"
	}
	protocol, ok := service.Annotations[AnnotationLbProtocol]
	if !ok {
		return "TCP"
	}
	return protocol
}

// 获取可用区
func (l *LoadBalancer) getAvailableZone(ctx context.Context) (string, string, error) {
	nodeList, err := l.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to list nodes: %w", err)
	}

	if nodeList == nil || len(nodeList.Items) == 0 {
		return "", "", errors.New("no nodes found in the cluster")
	}

	var (
		regionCode, azCode string
	)
	for _, node := range nodeList.Items {
		if node.Labels != nil {
			if az, ok := node.Labels[LabelNodeAzCode]; ok && az != "" {
				azCode = strings.TrimSpace(az)
			}
			if region, ok := node.Labels[LabelNodeRegionCode]; ok && region != "" {
				regionCode = strings.TrimSpace(region)
			}
		}
		if regionCode != "" && azCode != "" {
			break
		}
	}

	return regionCode, azCode, nil
}

// 获取计费方案ID
func (l *LoadBalancer) getBillingSchemeId(azCode, lbSpec string) (string, error) {
	req := lb.NewVpcSlbBillingSchemeRequest()
	req.AvailableZoneCode = azCode
	req.NetType = LbNetTypeWan
	req.BillingMethod = LbBillingMethodCostPay

	lsbSchema, err := api.VpcSlbBillingScheme(req)
	if err != nil {
		return "", fmt.Errorf("failed to get billing scheme: %w", err)
	}

	if lbSpec == "" {
		return "", fmt.Errorf("missing required annotation: %s", AnnotationLbSpec)
	}

	expectedSpecName, exists := lbSpecMap[lbSpec]
	if !exists {
		return "", fmt.Errorf("not fount lb spec: %s", lbSpec)
	}

	for _, schema := range lsbSchema.Data {
		if strings.Contains(schema.ConfName, expectedSpecName) {
			return schema.BillingSchemeId, nil
		}
	}

	return "", fmt.Errorf("billing scheme not found for spec: %s", expectedSpecName)
}

// 获取带宽计费方案ID
func (l *LoadBalancer) getBandwidthBillingSchemeId(azCode, billingMethod string, billingMethodConfId int) (string, error) {
	req := lb.NewBandwidthBillingSchemeRequest()
	req.AvailableZoneCode = azCode
	req.VpcId = consts.VpcID
	req.Type = BandwidthShared

	bandwidthResp, err := api.VpcBandwidthBillingScheme(req)
	if err != nil || bandwidthResp == nil {
		return "", fmt.Errorf("failed to query bandwidth billing scheme: %w", err)
	}

	if bandwidthResp.Code != consts.LbRequestSuccess {
		return "", fmt.Errorf("bandwidth billing scheme query failed with code: %s", bandwidthResp.Code)
	}

	if len(bandwidthResp.Data) == 0 {
		return "", errors.New("no bandwidth billing schemes available")
	}

	// 查找指定计费类型
	for _, bandwidth := range bandwidthResp.Data {
		if billingMethodConfId != 0 && billingMethodConfId != bandwidth.ConfId {
			continue
		}
		for _, bill := range bandwidth.BillingScheme {
			if bill.BillingType == billingMethod {
				return bill.BillingSchemeId, nil
			}
		}
	}
	return "", nil
}

// 构建创建SLB请求
func (l *LoadBalancer) buildCreateSlbRequest(service *v1.Service, params *serviceParams, azCode, billingSchemeId, bandwidthBillingSchemeId string) *lb.PackageCreateSlbRequest {
	slbName := SlbName(service.Name, service.Namespace, string(service.UID))
	netType := LbNetTypeWan
	if params.lbNetworkType == AllNetwork {
		netType = LbNetTypeWanLan
	}

	request := lb.NewPackageCreateSlbRequest()
	request.AvailableZoneCode = azCode
	request.VpcId = consts.VpcID
	request.Level = int(params.lbType)
	request.SlbInfo = lb.PackageCreateSlbInfo{
		BillingSchemeId: billingSchemeId,
		NetType:         netType,
		Name:            slbName,
		SubjectId:       params.subjectId,
	}

	request.BandwidthInfo = &lb.PackageCreateSlbBandwidthInfo{
		Name:            slbName,
		BillingSchemeId: bandwidthBillingSchemeId,
		Qos:             int(params.lbBandwidth),
		Type:            BandwidthShared,
		IsAutoRenewal:   false,
		IsToMonth:       false,
		Duration:        0,
		EipCount:        1,
		SubjectId:       params.subjectId,
	}

	return request
}

func (l *LoadBalancer) updateLbListen(ctx context.Context, service *v1.Service, nodes []*v1.Node, slbInfo *lb.DescribeVpcSlbResponseSlbInfo) ([]*lb.DescribeVpcSlbResponseVipInfo, error) {
	// 获取VIP地址
	params, err := l.parseServiceParams(service)
	if err != nil {
		return nil, err
	}

	vipList, err := l.extractVipFromResponse(slbInfo, params)
	if err != nil {
		return nil, err
	}

	if len(vipList) == 0 {
		return nil, errors.New("SLB ip resources are not ready")
	}

	// 获取调度算法
	algorithm := l.getSchedulerAlgorithm(service)

	var nodeInfo []string
	for _, nodeData := range nodes {
		nodeInfo = append(nodeInfo, nodeData.Name)
	}

	klog.Infof("update lb service listen, ns:%s, name:%s, externalTrafficPolicy:%v, target nodes:%s",
		service.Namespace, service.Name, service.Spec.ExternalTrafficPolicy, strings.Join(nodeInfo, ","))

	// 构建监听器列表
	listeners, err := l.buildListeners(service, nodes, algorithm, vipList)
	if err != nil {
		return nil, fmt.Errorf("failed to build listeners: %w", err)
	}

	if l.checkUpdateConforming(service, vipList, listeners) {
		klog.Info("listener conforming, skip update")
		return vipList, nil
	}

	operatorType := UpdateListenExact
	if params.selectSLB == "" {
		operatorType = UpdateListenFull
	}

	// 更新负载均衡监听器
	err = l.updateSlbListeners(slbInfo.SlbId, operatorType, listeners)
	if err != nil {
		return nil, err
	}

	return vipList, nil
}

func (l *LoadBalancer) checkUpdateConforming(service *v1.Service, vipList []*lb.DescribeVpcSlbResponseVipInfo, needChangeListeners []lb.VpcSlbUpdateListenRequestListen) (conforming bool) {
	conforming = true
	var needChangeListenersMap = make(map[string]lb.VpcSlbUpdateListenRequestListen)
	for _, listener := range needChangeListeners {
		needChangeListenersMap[listener.ListenIp] = listener
	}

	if len(vipList)*len(service.Spec.Ports) < len(needChangeListeners) {
		klog.Warningf("vipList length is less than needChangeListeners, vipList:%d, needChangeListeners:%d", len(vipList), len(needChangeListeners))
		conforming = false
		return conforming
	}

	for _, vipInfo := range vipList {
		listener, ok := needChangeListenersMap[vipInfo.Vip]
		if !ok {
			continue
		}
		conforming = l.checkListenerConforming(vipInfo.ListenList, listener)
		if !conforming {
			break
		}
	}
	return
}

func (l *LoadBalancer) checkListenerConforming(vipListenList []lb.ListenData, needChangeListener lb.VpcSlbUpdateListenRequestListen) (conforming bool) {
	conforming = false
	var targetListen *lb.ListenData
	for i, vipListen := range vipListenList {
		listenPort := GetListenPort(vipListen.ListenPort)
		if listenPort == 0 || vipListen.ListenName == "" {
			klog.Warningf("failed to get VIP %s listen port(%v->%v) or name(%s) from SLB details",
				needChangeListener.ListenIp, vipListen.ListenPort, listenPort, vipListen.ListenName)
			continue
		}
		if listenPort == needChangeListener.ListenPort && vipListen.ListenName == needChangeListener.ListenName {
			targetListen = &vipListenList[i]
			break
		}
	}
	if targetListen == nil {
		klog.Infof("need change listener %s is not found in vip listen list", needChangeListener.ListenName)
		return
	}
	// 首先比较长度
	if len(targetListen.RsList) != len(needChangeListener.RsList) {
		return
	}

	// 创建map方便比较
	actualRsMap := make(map[string]lb.DescribeVpcSlbRsInfo)
	for _, rs := range targetListen.RsList {
		actualRsMap[rs.RsIp] = rs
	}

	// 遍历期望的RS列表，检查是否都存在于实际RS列表中
	for _, needChangeRs := range needChangeListener.RsList {
		rs, exists := actualRsMap[needChangeRs.RsLanIp]
		if !exists {
			return
		}
		if rs.RsPort != fmt.Sprintf("%d", needChangeRs.RsPort) || rs.RsIp != needChangeRs.RsLanIp {
			return
		}

	}
	conforming = true

	klog.Infof("ip %s port %v listeners %s are conforming",
		needChangeListener.ListenIp, needChangeListener.ListenPort, needChangeListener.RsListString())

	return
}

// extractVipFromResponse 从响应中提取VIP地址
func (l *LoadBalancer) extractVipFromResponse(slbInfo *lb.DescribeVpcSlbResponseSlbInfo, params *serviceParams) ([]*lb.DescribeVpcSlbResponseVipInfo, error) {
	if slbInfo == nil || params == nil {
		return nil, errors.New("slbInfo cannot be nil")
	}

	// 根据是否指定特定SLB来筛选IP
	if params.selectSLB != "" {
		return l.filterSelectedVips(slbInfo.VipList, params)
	}

	// 根据网络类型进一步过滤IP
	return l.filterVipsByNetworkType(slbInfo.VipList, params)
}

// filterSelectedVips 根据用户指定的公网IP和私网IP筛选VIP
func (l *LoadBalancer) filterSelectedVips(vipList []lb.DescribeVpcSlbResponseVipInfo, params *serviceParams) ([]*lb.DescribeVpcSlbResponseVipInfo, error) {
	var selectedVips []lb.DescribeVpcSlbResponseVipInfo

	for _, vipInfo := range vipList {
		// 匹配公网IP (eip)
		if vipInfo.VipType == EIP && params.lbEip != "" && vipInfo.Vip == params.lbEip {
			selectedVips = append(selectedVips, vipInfo)
			continue
		}
		// 匹配私网IP (private)
		if vipInfo.VipType == LanVip && params.lbVip != "" && vipInfo.Vip == params.lbVip {
			selectedVips = append(selectedVips, vipInfo)
		}
	}

	return l.filterVipsByNetworkType(selectedVips, params)
}

// filterVipsByNetworkType 根据网络类型过滤VIP
func (l *LoadBalancer) filterVipsByNetworkType(vipList []lb.DescribeVpcSlbResponseVipInfo, params *serviceParams) ([]*lb.DescribeVpcSlbResponseVipInfo, error) {
	var public, private *lb.DescribeVpcSlbResponseVipInfo

	for i, vipInfo := range vipList {
		if len(params.ingressStatusIpMap) > 0 {
			_, ok := params.ingressStatusIpMap[vipInfo.Vip]
			if !ok {
				continue
			}
		}

		if vipInfo.VipType == EIP && public == nil {
			public = &vipList[i]
			continue
		}
		if vipInfo.VipType == LanVip && private == nil {
			private = &vipList[i]
			continue
		}

		if public != nil && private != nil {
			break
		}
	}

	switch params.lbNetworkType {
	case AllNetwork:
		if public != nil && private != nil {
			klog.Infof("public vip  %s, private vip  %s", public.Vip, private.Vip)
			return []*lb.DescribeVpcSlbResponseVipInfo{public, private}, nil
		}

	case PublicNetwork:
		if public != nil {
			klog.Infof("public vip  %s", public.Vip)
			return []*lb.DescribeVpcSlbResponseVipInfo{public}, nil
		}

	case PrivateNetwork:
		if private != nil {
			klog.Infof("private vip  %s", private.Vip)
			return []*lb.DescribeVpcSlbResponseVipInfo{private}, nil
		}

	default:
		return nil, fmt.Errorf("invalid network type: %s", params.lbNetworkType)
	}

	return nil, nil
}

// getSchedulerAlgorithm 获取调度算法，如果没有设置则使用默认值
func (l *LoadBalancer) getSchedulerAlgorithm(service *v1.Service) string {
	algorithm, ok := service.Annotations[AnnotationLbAlgorithm]
	if !ok {
		return LbAlgorithmRr
	}
	switch algorithm {
	case LbAlgorithmRr, LbAlgorithmWrr, LbAlgorithmHash:
		return algorithm
	default:
		return LbAlgorithmRr
	}
}

// buildListeners 构建监听器列表
func (l *LoadBalancer) buildListeners(service *v1.Service, nodes []*v1.Node, algorithm string, vipList []*lb.DescribeVpcSlbResponseVipInfo) ([]lb.VpcSlbUpdateListenRequestListen, error) {
	var listeners []lb.VpcSlbUpdateListenRequestListen
	protocol := l.getProtocol(service)
	for _, port := range service.Spec.Ports {
		for _, vip := range vipList {
			listener, err := l.buildListener(&port, nodes, protocol, algorithm, vip)
			if err != nil {
				return nil, fmt.Errorf("failed to build listener for port %d: %w", port.Port, err)
			}
			listeners = append(listeners, listener)
		}
	}

	return listeners, nil
}

// buildListener 构建单个监听器
func (l *LoadBalancer) buildListener(port *v1.ServicePort, nodes []*v1.Node, protocol, algorithm string, vip *lb.DescribeVpcSlbResponseVipInfo) (lb.VpcSlbUpdateListenRequestListen, error) {
	// 生成监听器名称
	listenName := l.generateListenName(vip.Vip, port.Port)

	// 构建真实服务器列表
	rsList, err := l.buildRealServerList(nodes, port)
	if err != nil {
		return lb.VpcSlbUpdateListenRequestListen{}, err
	}

	// 创建监听器对象
	listener := lb.VpcSlbUpdateListenRequestListen{
		ListenIp:       vip.Vip,
		ListenPort:     int(port.Port),
		ListenProtocol: protocol,
		Scheduler:      algorithm,
		ListenName:     listenName,
		Timeout:        10, // 默认超时时间10秒
		RsList:         rsList,
		HealthCheck: lb.VpcSlbUpdateListenRequestHealthCheck{
			Protocol:         protocol,
			ConnectTimeout:   5,
			Retry:            3,
			DelayLoop:        10,
			DelayBeforeRetry: 30,
		},
	}

	return listener, nil
}

// generateListenName 生成监听器名称
func (l *LoadBalancer) generateListenName(vip string, port int32) string {
	// 移除VIP中的点号，避免名称中出现特殊字符
	nameBase := fmt.Sprintf("eks-lb%s-%v", strings.ReplaceAll(vip, ".", ""), port)

	// 确保名称不超过25个字符的限制
	if len(nameBase) > 25 {
		return nameBase[:25]
	}

	return nameBase
}

// buildRealServerList 构建真实服务器列表
func (l *LoadBalancer) buildRealServerList(nodes []*v1.Node, port *v1.ServicePort) ([]lb.VpcSlbUpdateListenRequestRs, error) {
	var rsList []lb.VpcSlbUpdateListenRequestRs

	for _, node := range nodes {
		address, err := l.getNodeInternalAddress(node)
		if err != nil {
			return nil, err
		}

		realServer := lb.VpcSlbUpdateListenRequestRs{
			RsId:    node.Spec.ProviderID,
			RsName:  node.Name,
			RsType:  EKS,
			RsLanIp: address,
			RsPort:  int(port.NodePort),
			Weight:  50, // 默认权重为50，后续可考虑根据节点上的Pod数量动态调整
		}

		rsList = append(rsList, realServer)
	}

	return rsList, nil
}

// getNodeInternalAddress 获取节点的内部IP地址
func (l *LoadBalancer) getNodeInternalAddress(node *v1.Node) (string, error) {
	if node.Spec.ProviderID == "" {
		return "", fmt.Errorf("node %s does not have a provider ID", node.Name)
	}
	for _, addr := range node.Status.Addresses {
		if addr.Type == IpTypeInternal && addr.Address != "" {
			return addr.Address, nil
		}
	}

	return "", fmt.Errorf("node %s does not have an internal IP address", node.Name)
}

// updateSlbListeners 更新负载均衡器的监听器
func (l *LoadBalancer) updateSlbListeners(slbId, operatorType string, listeners []lb.VpcSlbUpdateListenRequestListen) error {
	// 如果没有监听器需要更新，直接返回
	if len(listeners) == 0 {
		klog.Infof("No listeners to update for SLB ID: %s", slbId)
		return nil
	}

	// 创建更新请求
	request := lb.NewVpcSlbUpdateListenRequest()
	request.ListenList = listeners
	request.SlbId = slbId
	request.Platform = EKS
	request.OperatorType = operatorType

	// 执行API调用
	response, err := api.VpcSlbUpdateListen(request)
	if err != nil {
		return fmt.Errorf("failed to update SLB listeners: %w", err)
	}

	// 等待任务完成
	if err = l.describeTask(response.TaskId); err != nil {
		return errors.New("waiting for SLB listener update task to complete")
	}

	klog.Infof("Successfully updated %d listeners for SLB ID: %s", len(listeners), slbId)
	return nil
}

func (l *LoadBalancer) clearLbListen(ctx context.Context, service *v1.Service, slbInfo *lb.DescribeVpcSlbResponseSlbInfo) error {
	if len(slbInfo.VipList) == 0 {
		klog.Infof("No listeners to clear for SLB ID: %s, skip.", slbInfo.SlbId)
		return nil
	}
	var (
		err              error
		listenIds        []string
		taskId           string
		deleteListenResp *lb.DeleteVpcSLBListenResponse
		clearResp        *lb.VpcSlbClearListenResponse
	)

	listenIds, err = l.getSelfListen(ctx, service, slbInfo)
	if len(listenIds) != 0 {
		req := lb.NewDeleteLbListenersRequest()
		req.ListenIds = listenIds
		deleteListenResp, err = api.DeleteVpcSLBListenRequest(req)
		if err != nil {
			return err
		}
		taskId = deleteListenResp.TaskId
	} else {
		clearResp, err = api.VpcSlbClearListen(slbInfo.SlbId)
		if err != nil {
			if clearResp != nil && clearResp.Code == consts.ErrorSlbNotFound {
				return nil
			}
			return err
		}
		taskId = clearResp.TaskId
	}

	return l.describeTask(taskId)
}

func (l *LoadBalancer) describeLbInstance(ctx context.Context, service *v1.Service) (*lb.DescribeVpcSlbResponse, error) {
	request := lb.NewDescribeVpcSlbRequest()
	if service.Annotations != nil && service.Annotations[AnnotationLbId] != "" {
		request.SlbID = service.Annotations[AnnotationLbId]
	} else {
		request.SlbName = SlbName(service.Name, service.Namespace, string(service.UID))
	}
	response, err := api.DescribeVpcSlb(request)
	if response != nil && response.Code == consts.ErrorSlbNotFound {
		klog.Warningf("%v, id:%s, name:%s", SLBNotFound, request.SlbID, request.SlbName)
		return nil, SLBNotFound
	}
	if err != nil {
		return nil, err
	}
	// 接口请求返回异常
	if response != nil && response.Code != consts.LbRequestSuccess {
		klog.Errorf("DescribeVpcSlb failed, msg: %s", response.Message)
		return nil, errors.New(response.Message)
	}
	return response, nil
}

func (l *LoadBalancer) getSelfListen(ctx context.Context, service *v1.Service, slbInfo *lb.DescribeVpcSlbResponseSlbInfo) ([]string, error) {
	var (
		ipMap     = make(map[string]bool)
		listenIds []string
	)

	if len(service.Status.LoadBalancer.Ingress) == 0 {
		return nil, nil
	}

	for _, ingress := range service.Status.LoadBalancer.Ingress {
		if ingress.IP == "" {
			continue
		}
		ipMap[ingress.IP] = true
	}

	for _, v := range slbInfo.VipList {
		_, ok := ipMap[v.Vip]
		if !ok {
			continue
		}
		headName := "eks-lb" + strings.ReplaceAll(v.Vip, ".", "")
		for _, listen := range v.ListenList {
			if listen.ListenName == "" || listen.ListenId == "" {
				continue
			}
			if strings.Contains(listen.ListenName, headName) {
				listenIds = append(listenIds, listen.ListenId)
			}
		}
	}
	return listenIds, nil
}

func (l *LoadBalancer) describeTask(taskId string) error {
	if taskId == "" {
		return errors.New("taskId is empty")
	}
	for i := 0; i < 200; i++ {
		resp, err := api.DescribeTask(taskId)
		if err != nil {
			klog.Warningf("describe task failed, error: %v", err)
			time.Sleep(time.Second * 3)
			continue
		}
		switch resp.Data.TaskStatus {
		case LbTaskSuccess:
			return nil
		case LbTakError:
			return errors.New("SLB resource creation failed")
		default:
			time.Sleep(time.Second * 3)
		}
	}
	return errors.New("SLB resource creation in progress")
}

// SlbName 通过hash值的方式，计算slb的名称，让slb名称具有一致性和独立性
func SlbName(svcName, namespace, uid string) string {
	hash := sha256.New()
	hash.Write([]byte(strings.Trim(string(uid), "-")))
	value := hash.Sum(nil)
	name := fmt.Sprintf("%s-%s-%s", svcName, namespace, hex.EncodeToString(value)[:16])
	if len(name) > 64 {
		name = name[len(name)-64:]
	}
	return name
}

func GetListenPort(port interface{}) int {
	// SLB OPEN-API ListenPort存在两个版本不一致的数据类型，需适配断言处理
	switch port.(type) {
	case int:
		return port.(int)
	case int64:
		return int(port.(int64))
	default:
		portStr := fmt.Sprintf("%v", port)
		klog.Warningf("port: %v->%v", port, portStr)
		portInt64, _ := strconv.ParseInt(portStr, 10, 64)
		return int(portInt64)
	}
}
