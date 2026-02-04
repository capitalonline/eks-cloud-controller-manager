package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/controller"
	_ "github.com/capitalonline/eks-cloud-controller-manager/pkg/provider"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	"k8s.io/cloud-provider/options"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version"
	"k8s.io/klog/v2"
)

func main() {
	klog.Info("ccm start...")
	rand.Seed(time.Now().UTC().UnixNano())
	logs.InitLogs()
	defer logs.FlushLogs()
	opts, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}
	controllerInitializers := initFuncConstructors()

	fss := cliflag.NamedFlagSets{}
	fss.FlagSet(consts.ProviderName)
	command := app.NewCloudControllerManagerCommand(opts, cloudInitializer, controllerInitializers, fss, wait.NeverStop)

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	if err = command.Execute(); err != nil {
		klog.Fatalf("unable to execute command: %v", err)
	}
}

func cloudInitializer(config *cloudcontrollerconfig.CompletedConfig) cloudprovider.Interface {
	cloudConfig := config.ComponentConfig.KubeCloudShared.CloudProvider
	klog.Info("cloudConfig ", cloudConfig)
	providerName := cloudConfig.Name
	if providerName == "" {
		providerName = consts.ProviderName
	}
	cloud, err := cloudprovider.InitCloudProvider(consts.ProviderName, cloudConfig.CloudConfigFile)

	if err != nil {
		klog.Fatalf("Cloud provider could not be initialized: %v", err)
	}
	if cloud == nil {
		klog.Fatalf("Cloud provider is nil")
	}

	if !cloud.HasClusterID() {
		if config.ComponentConfig.KubeCloudShared.AllowUntaggedCloud {
			klog.Warning("detected a cluster without a ClusterID.  A ClusterID will be required in the future.  Please tag your cluster to avoid any future issues")
		} else {
			klog.Fatalf("no ClusterID found.  A ClusterID is required for the cloud provider to function properly.  This check can be bypassed by setting the allow-untagged-cloud option")
		}
	}

	return cloud
}

func initFuncConstructors() map[string]app.ControllerInitFuncConstructor {
	defaultInitFuncConstructors := app.DefaultInitFuncConstructors
	nodeControllerWrapper := new(controller.NodeControllerWrapper)
	// 注册一个新的节点同步器
	defaultInitFuncConstructors[controller.NodeControllerKey] = app.ControllerInitFuncConstructor{
		InitContext: app.ControllerInitContext{
			ClientName: controller.NodeControllerClientName,
		},
		Constructor: nodeControllerWrapper.StartNodeControllerWrapper,
	}
	return defaultInitFuncConstructors

}
