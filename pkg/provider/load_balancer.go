package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/api"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/lb"
	v1 "k8s.io/api/core/v1"
	"math/rand"
	"strings"
	"time"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"crypto/sha256"
	"encoding/hex"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"strconv"
)

const (
	AnnotationLbProtocol      = "service.beta.kubernetes.io/cds-load-balancer-protocol"
	AnnotationLbType          = "service.beta.kubernetes.io/cds-load-balancer-types"
	AnnotationLbSpec          = "service.beta.kubernetes.io/cds-load-balancer-specification"
	AnnotationLbBandwidth     = "service.beta.kubernetes.io/cds-load-balancer-bandwidth"
	AnnotationLbEip           = "service.beta.kubernetes.io/cds-load-balancer-eip"
	AnnotationLbAlgorithm     = "service.beta.kubernetes.io/cds-load-balancer-algorithm"
	AnnotationLbSubjectId     = "service.beta.kubernetes.io/cds-load-balancer-subject-id"
	AnnotationLbBillingMethod = "service.beta.kubernetes.io/cds-load-balancer-billingmethod"
	AnnotationLbListen        = "service.eks.listen"
	LbNetTypeWan              = "wan"
	LbNetTypeWanLan           = "wan_lan"
	LbBillingMethodCostPay    = "0" // 按需计费
	LabelNodeAz               = "node.kubernetes.io/node.az"
	LabelNodeAzCode           = "node.kubernetes.io/node.az-code"
	LbTaskSuccess             = "success"
	LbTakError                = "error"
	DefaultBillingType        = "number"
	BandwidthShared           = "shared"
)

const (
	LBSpecStandard = "standard" // 标准型
	LBSpecHigh     = "high"     // 高阶型
	LBSpecSuper    = "super"    // 超强型
	LBSpecExtreme  = "extreme"  // 至强型

	//LBSpecNameStandard = "标准型Ⅰ"
	//LBSpecNameHigh     = "高阶型Ⅰ"
	//LBSpecNameSuper    = "超强型Ⅰ"
	//LBSpecNameExtreme  = "至强型Ⅰ"

	LBSpecNameStandard = "标准型"
	LBSpecNameHigh     = "高阶型"
	LBSpecNameSuper    = "超强型"
	LBSpecNameExtreme  = "至强型"
)

const (
	IpTypeInternal = "InternalIP"
	PlatformEks    = "eks"

	UpdateListenFull  = "full"  // 全量更新
	UpdateListenExact = "exact" // 精确更新

	RsTypeEks = "eks"
)

var lbSpecMap = map[string]string{
	LBSpecStandard: LBSpecNameStandard,
	LBSpecHigh:     LBSpecNameHigh,
	LBSpecSuper:    LBSpecNameSuper,
	LBSpecExtreme:  LBSpecNameExtreme,
}

type LoadBalancer struct {
	clientSet *kubernetes.Clientset
}

// GetLoadBalancer 查询lb
func (l *LoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {

	response, err := l.describeLbInstance(ctx, clusterName, service)
	// k8s在删除节点之后会查一遍slb，确认是否被删除
	if response != nil && response.Code == "50002" {
		return nil, false, nil
	}
	if err != nil {
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
	response, err := l.describeLbInstance(ctx, clusterName, service)
	if err != nil {
		return ""
	}
	slb := response.Data
	return slb.SlbName
}

// EnsureLoadBalancer 创建lb
func (l *LoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {

	if service.Spec.SessionAffinity != v1.ServiceAffinityNone {
		return nil, errors.New("SessionAffinity is not supported currently, only support 'None' type")
	}

	// 先查询lb是不是已经创建过了
	request := lb.NewDescribeVpcSlbRequest()
	request.SlbName = SlbName(service.Name, service.Namespace, string(service.UID))
	response, err := api.DescribeVpcSlb(request)
	// 调用接口有问题，接口没返回json
	if err != nil && response == nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("查询slb失败")
	}
	// 接口请求返回异常
	if response.Code != consts.LbRequestSuccess && response.Code != consts.ErrorSlbNotFound {
		klog.Error(fmt.Sprintf("查询slb失败，response: %#v", response))
		return nil, errors.New(response.Message)
	}
	var slbId string
	// lb不存在
	if len(response.Data.SlbId) < 1 {
		slbId, err = l.createSlb(ctx, service, nodes)
		if err != nil {
			return nil, err
		}
	} else {
		slbId = response.Data.SlbId
	}
	err = l.UpdateLoadBalancer(ctx, clusterName, service, nodes)
	if err != nil {
		return nil, err
	}
	// 重新查一遍slb信息
	request = lb.NewDescribeVpcSlbRequest()
	request.SlbName = SlbName(service.Name, service.Namespace, string(service.UID))
	request.SlbID = slbId
	response, err = api.DescribeVpcSlb(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("查询slb异常")
	}
	ingresses := make([]v1.LoadBalancerIngress, 0, len(response.Data.VipList))
	for i := 0; i < len(response.Data.VipList); i++ {
		vip := response.Data.VipList[i]
		ingresses = append(ingresses, v1.LoadBalancerIngress{IP: vip.Vip})
	}
	return &v1.LoadBalancerStatus{
		Ingress: ingresses,
	}, nil
}

func (l *LoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return l.updateLbListen(ctx, service, nodes)
}

func (l *LoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	return l.clearLbListen(ctx, clusterName, service)
}

// 创建slb
func (l *LoadBalancer) createSlb(ctx context.Context, service *v1.Service, nodes []*v1.Node) (string, error) {
	if len(service.Annotations) == 0 {
		return "", errors.New("service annotations is null")
	}
	var subjectId int
	// 测试金查询
	if service.Annotations[AnnotationLbSubjectId] != "" {
		i, _ := strconv.ParseInt(service.Annotations[AnnotationLbSubjectId], 10, 64)
		subjectId = int(i)
	}
	// 协议
	//proctol := service.Annotations[AnnotationLbProtocol]
	lbType, err := strconv.ParseInt(service.Annotations[AnnotationLbType], 10, 64)
	if err != nil {
		return "", err
	}
	lbSpec := service.Annotations[AnnotationLbSpec]
	//lbBandwidth := service.Annotations[AnnotationLbBandwidth]

	billingMethod := service.Annotations[AnnotationLbBillingMethod]
	if billingMethod == "" {
		billingMethod = DefaultBillingType
	}

	lbBandwidth, err := strconv.ParseInt(service.Annotations[AnnotationLbBandwidth], 10, 64)
	if err != nil {
		return "", err
	}

	lbEip, err := strconv.ParseInt(service.Annotations[AnnotationLbEip], 10, 64)
	if err != nil {
		return "", err
	}

	// 随机获取一个节点所在可用区
	nodeList, err := l.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if nodeList == nil || len(nodeList.Items) == 0 {
		return "", errors.New("查询节点信息异常，请稍后重试")
	}
	var azList = make([]string, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		if len(node.Labels) == 0 {
			continue
		}
		if az, ok := node.Labels[LabelNodeAzCode]; ok {
			azList = append(azList, az)
		}
	}
	if len(azList) == 0 {
		return "", errors.New("获取az信息失败")
	}
	randomInt := rand.Intn(len(azList))
	var azCode = azList[randomInt]
	azCode = strings.TrimSpace(azCode)
	// 查询计费方案
	lsbSchemaReq := lb.NewVpcSlbBillingSchemeRequest()
	lsbSchemaReq.AvailableZoneCode = azCode
	lsbSchemaReq.NetType = LbNetTypeWan
	lsbSchemaReq.BillingMethod = LbBillingMethodCostPay
	lsbSchema, err := api.VpcSlbBillingScheme(lsbSchemaReq)
	if err != nil {
		return "", err
	}
	var billingSchemeId string

	for i := 0; i < len(lsbSchema.Data); i++ {
		schema := lsbSchema.Data[i]
		if strings.Contains(schema.ConfName, lbSpecMap[lbSpec]) {
			billingSchemeId = schema.BillingSchemeId
			break
		}
	}
	if billingSchemeId == "" {
		return "", errors.New("未查到相关计费信息")
	}

	// 轮训策略
	//lbAlgorithm := service.Annotations[AnnotationLbAlgorithm]
	request := lb.NewPackageCreateSlbRequest()
	// 获取一个azcode
	request.AvailableZoneCode = azCode
	request.VpcId = consts.VpcID
	// 当前仅支持4层
	request.Level = int(lbType)
	request.SlbInfo = lb.PackageCreateSlbInfo{
		BillingSchemeId: billingSchemeId,
		NetType:         LbNetTypeWan,
		Name:            SlbName(service.Name, service.Namespace, string(service.UID)),
		SubjectId:       subjectId,
	}
	// 查询共享带宽计费ID
	bandwidthReq := lb.NewBandwidthBillingSchemeRequest()
	// 获取RegionCode
	bandwidthReq.AvailableZoneCode = azCode
	bandwidthReq.VpcId = consts.VpcID
	bandwidthReq.Type = BandwidthShared
	bandwidthResp, err := api.VpcBandwidthBillingScheme(bandwidthReq)
	if err != nil || bandwidthResp == nil {
		return "", fmt.Errorf("查询共享带宽计费失败,err:%v", err)
	}
	if bandwidthResp.Code != consts.LbRequestSuccess {
		return "", errors.New(fmt.Sprintf("查询共享带宽计费失败，code:%s", bandwidthResp.Code))
	}
	bandwidthBillingSchemeId := ""

outer:
	for i := 0; i < len(bandwidthResp.Data); i++ {
		bandwidth := bandwidthResp.Data[i]
		for j := 0; j < len(bandwidth.BillingScheme); j++ {
			bill := bandwidth.BillingScheme[j]
			if bill.BillingType == billingMethod {
				bandwidthBillingSchemeId = bill.BillingSchemeId
				break outer
			}
		}
	}
	//	 如果没有number类型计费，默认拿第一种方案的计费
	if bandwidthBillingSchemeId == "" {
		if len(bandwidthResp.Data) == 0 || len(bandwidthResp.Data[0].BillingScheme) == 0 {
			return "", fmt.Errorf("没有相关的共享带宽计费方案")
		}
		bandwidthBillingSchemeId = bandwidthResp.Data[0].BillingScheme[0].BillingSchemeId
	}

	request.BandwidthInfo = lb.PackageCreateSlbBandwidthInfo{
		Name:            SlbName(service.Name, service.Namespace, string(service.UID)),
		BillingSchemeId: bandwidthBillingSchemeId,
		Qos:             int(lbBandwidth),
		Type:            BandwidthShared,
		IsAutoRenewal:   false,
		IsToMonth:       false,
		Duration:        0,
		EipCount:        int(lbEip),
		SubjectId:       subjectId,
	}

	response, err := api.PackageCreateSlb(request)
	if err != nil || response == nil {
		return "", fmt.Errorf("创建slb失败:%v", err)
	}
	if response.Code != consts.LbRequestSuccess {
		return "", errors.New(fmt.Sprintf("创建lb失败%v %v", err, response.Message))
	}
	return response.Data.SlbId, l.describeTask(response.TaskId)
}

func (l *LoadBalancer) updateLbListen(ctx context.Context, service *v1.Service, nodes []*v1.Node) error {
	//listeners := make([]lb.VpcSlbUpdateListenRequestListen, 0, len(service.Spec.Ports))
	// 查询service的Annotations是否保存有上次更改时的记录
	//var listenList = make([]lb.VpcSlbUpdateListenRequestListen, 0, len(service.Spec.Ports))
	//var newListenList = make([]lb.VpcSlbUpdateListenRequestListen, 0, len(service.Spec.Ports)) // 用来存储新的值
	listenAnnotation, ok := service.Annotations[AnnotationLbListen]
	if ok {
		if err := json.Unmarshal([]byte(listenAnnotation), &listenAnnotation); err != nil {
			return err
		}
	}
	algorithm, ok := service.Annotations[AnnotationLbAlgorithm]
	if !ok {
		algorithm = "rr"
	}

	lbResp, err := l.describeLbInstance(ctx, "", service)
	if err != nil {
		return err
	}

	var vip string
	for i := 0; i < len(lbResp.Data.VipList); i++ {
		if lbResp.Data.VipList[i].Vip != "" {
			vip = lbResp.Data.VipList[i].Vip
			break
		}
	}
	if vip == "" {
		return errors.New("未查询到slb的vip")
	}
	var createList = make([]lb.VpcSlbUpdateListenRequestListen, 0, len(service.Spec.Ports))
	var updateList = make([]lb.VpcSlbUpdateListenRequestListen, 0, len(service.Spec.Ports))
	for i := 0; i < len(service.Spec.Ports); i++ {
		port := service.Spec.Ports[i]

		listenName := fmt.Sprintf("勿删-%s-%v", strings.Replace(vip, ".", "", -1), port.Port)
		if len(listenName) > 25 {
			// 接口限制监听名最长只能是25个字符
			listenName = listenName[:25]
		}
		listen := lb.VpcSlbUpdateListenRequestListen{
			ListenIp:       vip,
			ListenPort:     int(port.Port),
			ListenProtocol: string(port.Protocol),
			Scheduler:      algorithm,
			ListenName:     listenName,
			Timeout:        10, // 默认超时时间10
			RsList:         nil,
		}
		rsList := make([]lb.VpcSlbUpdateListenRequestRs, 0, len(nodes))
		for j := 0; j < len(nodes); j++ {
			node := nodes[j]
			var address string
			for k := 0; k < len(node.Status.Addresses); k++ {
				addr := node.Status.Addresses[k]
				if addr.Type == IpTypeInternal {
					address = addr.Address
					break
				}
			}
			rsList = append(rsList, lb.VpcSlbUpdateListenRequestRs{
				RsId:    node.Spec.ProviderID,
				RsName:  node.Name,
				RsType:  RsTypeEks,
				RsLanIp: address,
				RsPort:  int(port.NodePort),
				// TODO 根据接点上pod数量给权重
				Weight: 50,
			})

		}
		listen.RsList = rsList

		listen.HealthCheck = lb.VpcSlbUpdateListenRequestHealthCheck{
			Protocol:         string(port.Protocol),
			ConnectTimeout:   5,
			Retry:            3,
			DelayLoop:        10,
			DelayBeforeRetry: 30,
		}
		//if _, ok := listenList[string(port.NodePort)]; !ok {
		//updateList = append(createList, listen)
		//} else {
		updateList = append(updateList, listen)
		//}
		//newlistenInfo[string(port.NodePort)] = ""
	}

	if len(createList) > 0 {

		request := lb.NewVpcSlbUpdateListenRequest()
		request.ListenList = createList
		request.SlbId = lbResp.Data.SlbId
		request.Platform = PlatformEks
		request.OperatorType = UpdateListenExact
		response, err := api.VpcSlbUpdateListen(request)
		if err != nil {
			return err
		}
		return l.describeTask(response.TaskId)
	}
	if len(updateList) > 0 {
		request := lb.NewVpcSlbUpdateListenRequest()
		request.ListenList = updateList
		request.SlbId = lbResp.Data.SlbId
		request.Platform = PlatformEks
		request.OperatorType = UpdateListenFull
		response, err := api.VpcSlbUpdateListen(request)
		if err != nil {
			return err
		}
		return l.describeTask(response.TaskId)
	}
	// TODO 后续优化，精确更新更新slb
	return nil
}

func (l *LoadBalancer) clearLbListen(ctx context.Context, clusterName string, service *v1.Service) error {
	request := lb.NewDescribeVpcSlbRequest()
	request.SlbName = SlbName(service.Name, service.Namespace, string(service.UID))
	response, err := api.DescribeVpcSlb(request)
	//klog.Info(fmt.Sprintf("清除监听：%#v ,%v", response, err))
	if err != nil {
		klog.Errorf("清除监听失败,查询slb")
		return err
	}
	slb := response.Data
	clearResp, err := api.VpcSlbClearListen(slb.SlbId)
	// slb不存在，直接返回成功
	if clearResp != nil && clearResp.Code == "50002" {
		return nil
	}
	if err != nil {
		return err
	}
	return l.describeTask(clearResp.TaskId)
}

func (l *LoadBalancer) describeLbInstance(ctx context.Context, clusterName string, service *v1.Service) (*lb.DescribeVpcSlbResponse, error) {
	request := lb.NewDescribeVpcSlbRequest()
	request.SlbName = SlbName(service.Name, service.Namespace, string(service.UID))
	response, err := api.DescribeVpcSlb(request)
	if err != nil {
		return response, err
	}
	// 接口请求返回异常
	if response.Code != consts.LbRequestSuccess {
		klog.Error(fmt.Sprintf("查询slb失败，response: %#v", response))
		return nil, errors.New(response.Message)
	}
	return response, err
}

func (l *LoadBalancer) describeTask(taskId string) error {
	for i := 0; i < 600; i++ {
		resp, err := api.DescribeTask(taskId)
		if err != nil {
			return err
		}
		switch resp.Data.TaskStatus {
		case LbTaskSuccess:
			return nil
		case LbTakError:
			return errors.New("任务失败")
		default:
			time.Sleep(time.Second * 3)
		}
	}
	return errors.New("任务超时")
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
