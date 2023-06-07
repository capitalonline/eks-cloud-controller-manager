package main

import (
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/klog/v2"
	// For existing cloud providers, the option to import legacy providers is still available.
	// e.g. _"k8s.io/legacy-cloud-providers/<provider>"
)

func main() {
	//ccmOptions, err := options.NewCloudControllerManagerOptions()
	//if err != nil {
	//	klog.Fatalf("unable to initialize command options: %v", err)
	//}

	controllerInitializers := app.DefaultInitFuncConstructors

	//nodeIpamController := nodeIPAMController{}
	//nodeIpamController.nodeIPAMControllerOptions.NodeIPAMControllerConfiguration = &nodeIpamController.nodeIPAMControllerConfiguration
	//fss := cliflag.NamedFlagSets{}
	//nodeIpamController.nodeIPAMControllerOptions.AddFlags(fss.FlagSet("nodeipam controller"))

	controllerInitializers["nodeipam"] = app.ControllerInitFuncConstructor{
		// "node-controller" is the shared identity of all node controllers, including node, node lifecycle, and node ipam.
		// See https://github.com/kubernetes/kubernetes/pull/72764#issuecomment-453300990 for more context.
		InitContext: app.ControllerInitContext{
			ClientName: "node-controller",
		},
		//Constructor: nodeIpamController.StartNodeIpamControllerWrapper,
	}

	//command := app.NewCloudControllerManagerCommand(ccmOptions, cloudInitializer, controllerInitializers, fss, wait.NeverStop)
	//fmt.Println(command)
	//code := cli.Run(command)
	//os.Exit(code)
}

func cloudInitializer(config *cloudcontrollerconfig.CompletedConfig) cloudprovider.Interface {
	cloudConfig := config.ComponentConfig.KubeCloudShared.CloudProvider
	// initialize cloud provider with the cloud provider name and config file provided
	cloud, err := cloudprovider.InitCloudProvider(cloudConfig.Name, cloudConfig.CloudConfigFile)
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
