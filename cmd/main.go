package main

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/controller"

	_ "github.com/capitalonline/eks-cloud-controller-manager/pkg/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"math/rand"
	"time"
)

func main() {
	klog.Info("程序启动")
	rand.Seed(time.Now().UTC().UnixNano())
	logs.InitLogs()
	defer logs.FlushLogs()
	opts, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}
	opts.NodeStatusUpdateFrequency = metav1.Duration{Duration: time.Second * 30}
	controllerInitializers := app.DefaultInitFuncConstructors

	nodeController := controller.ControllerWrapper{}
	fss := cliflag.NamedFlagSets{}

	controllerInitializers[controller.NodeControllerKey] = app.ControllerInitFuncConstructor{
		InitContext: app.ControllerInitContext{
			ClientName: "node-controller",
		},
		Constructor: nodeController.StartNodeControllerWrapper,
	}
	fss.FlagSet(consts.ProviderName)
	//app.ControllersDisabledByDefault.Insert(controller.NodeControllerKey)
	command := app.NewCloudControllerManagerCommand(opts, cloudInitializer, controllerInitializers, fss, wait.NeverStop)
	command.Flags().Set("cloud-provider", "true")

	if err := command.Execute(); err != nil {
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
