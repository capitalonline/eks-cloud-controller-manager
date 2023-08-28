package consts

const ApiVersion = "2019-08-08"

const (
	ServiceEKS = "eks/v1"
	ServiceLb  = "vpc"
)
const ProviderName = "cdscloud"

const (
	ApiHostAddress   = "api.capitalonline.net"
	LbApiHostAddress = "cdsapi.capitalonline.net"
)

const (
	ActionDescribeEKSNode           = "DescribeEKSNode"     // 查询节点
	ActionNodeCCMInit               = "NodeCCMInit"         // 初始化ccm
	ActionModifyClusterLoad         = "ModifyClusterLoad"   // 更新节点负载
	ActionPackageCreateSlb          = "PackageCreateSlb"    // 创建负载均衡实例
	ActionDescribeLbInstance        = "DescribeVpcSlb"      // 查询负载均衡实例
	ActionVpcSlbUpdateListen        = "VpcSlbUpdateListen"  // 更新负载均衡实例
	ActionDeleteLbInstance          = "DeleteLbInstance"    // 删除负载均衡实例
	ActionDescribeLbTaskStatus      = "DescribeTask"        // 查询任务状态
	ActionVpcSlbClearListen         = "VpcSlbClearListen"   // 清空监听
	ActionVpcSlbBillingScheme       = "VpcSlbBillingScheme" // slb计费查询
	ActionVpcBandwidthBillingScheme = "BandwidthBillingScheme"
)

const (
	EnvAccessKeyID     = "CDS_ACCESS_KEY_ID"
	EnvAccessKeySecret = "CDS_ACCESS_KEY_SECRET"
	EnvRegion          = "CDS_REGION"
	EnvAPIHost         = "CDS_API_HOST"
	EnvClusterId       = "CDS_CLUSTER_ID"
	EnvAz              = "CDS_CLUSTER_AZ"
	SCHEMA             = "CDS_API_SCHEMA"
	EnvLbApiHost       = "CDS_LB_API_HOST"
	EnvVpcId           = "CDS_VPC_ID"
)

const (
	LbRequestSuccess = "Success"
	ErrorSlbNotFound = "50002"
)

const (
	NodeStatusError   = "error"
	NodeStatusDeleted = "deleted"
	NodeStatusFailed  = "failed"
	NodeStatusRunning = "running"
)
