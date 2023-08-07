package consts

const ApiVersion = "2019-08-08"

const (
	ServiceEKS = "eks/v1"
	ServiceLb  = "lb/v1"
)
const ProviderName = "cdscloud"

const (
	ApiHostAddress = "api.capitalonline.net"
)

const (
	ActionDescribeEKSNode      = "DescribeEKSNode"    // 查询节点
	ActionNodeCCMInit          = "NodeCCMInit"        // 初始化ccm
	ActionModifyClusterLoad    = "ModifyClusterLoad"  // 更新节点负载
	ActionCreateLb             = "CreateLb"           // 创建负载均衡实例
	ActionDescribeLbInstance   = "DescribeLbInstance" // 查询负载均衡实例
	ActionUpdateLbInstance     = "UpdateLbInstance"   // 更新负载均衡实例
	ActionDeleteLbInstance     = "DeleteLbInstance"   // 删除负载均衡实例
	ActionDescribeLbTaskStatus = "DescribeTaskStatus" // 查询任务状态
)

const (
	EnvAccessKeyID     = "CDS_ACCESS_KEY_ID"
	EnvAccessKeySecret = "CDS_ACCESS_KEY_SECRET"
	EnvRegion          = "CDS_REGION"
	EnvAPIHost         = "CDS_API_HOST"
	EnvClusterId       = "CDS_CLUSTER_ID"
	EnvAz              = "CDS_CLUSTER_AZ"
	SCHEMA             = "CDS_API_SCHEMA"
)
