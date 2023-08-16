package consts

import (
	"k8s.io/klog/v2"
	"os"
)

var (
	AccessKeyID     string = "ef646cc6132011eeaff59a519b3265fd"
	AccessKeySecret string = "0700d013dd2e4dd96d975de9694c0ccb"
	Region          string
	Az              string
	ClusterId       string
	ApiHost         string
	LbApiHost       string
	VpcID           string
)

func init() {
	ak := os.Getenv(EnvAccessKeyID)
	if ak == "" {
		klog.Infoln("未获取到ak")
		//panic("env CDS_ACCESS_KEY_ID must be set")
	}
	AccessKeyID = ak
	sk := os.Getenv(EnvAccessKeySecret)
	if sk == "" {
		klog.Infoln("未获取sk")
		//panic("env CDS_ACCESS_KEY_SECRET must be set")
	}
	AccessKeySecret = sk

	ClusterId = os.Getenv(EnvClusterId)

	if ClusterId == "" {
		klog.Infoln("未获取到集群id")
		//panic("env CDS_CLUSTER_ID must be set")
	}
	region := os.Getenv(EnvRegion)
	if region != "" {
		Region = region
	}
	az := os.Getenv(EnvAz)
	if az != "" {
		Az = az
	}
	apiHost := os.Getenv(EnvAPIHost)
	if apiHost != "" {
		ApiHost = apiHost
	} else {
		ApiHost = ApiHostAddress
	}
	if os.Getenv(EnvLbApiHost) != "" {
		LbApiHost = os.Getenv(EnvLbApiHost)
	} else {
		LbApiHost = LbApiHostAddress
	}
	if os.Getenv(EnvVpcId) != "" {
		VpcID = os.Getenv(EnvVpcId)
	}
}
