package consts

import "os"

var (
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	Az              string
	ClusterId       string
	ApiHost         string
)

func init() {
	ak := os.Getenv(EnvAccessKeyID)
	if ak == "" {
		panic("env CDS_ACCESS_KEY_ID must be set")
	}
	AccessKeyID = ak
	sk := os.Getenv(EnvAccessKeySecret)
	if sk == "" {
		panic("env CDS_ACCESS_KEY_SECRET must be set")
	}
	AccessKeySecret = sk

	clusterId := os.Getenv(EnvClusterId)
	if clusterId == "" {
		panic("env CDS_CLUSTER_ID must be set")
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
}
