package controller

import (
	"context"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	genericcontrollermanager "k8s.io/controller-manager/app"
	"k8s.io/controller-manager/controller"
)

const (
	NodeControllerClientName = "node-controller"

	NodeControllerKey = "label-taint"
)

type ControllerWrapper struct {
	//Options options.TaggingControllerOptions
}

func (c *ControllerWrapper) StartNodeControllerWrapper(initContext app.ControllerInitContext, completedConfig *cloudcontrollerconfig.CompletedConfig, cloud cloudprovider.Interface) app.InitFunc {
	return func(ctx context.Context, controllerContext genericcontrollermanager.ControllerContext) (controller.Interface, bool, error) {
		return c.startNodeController(ctx, initContext, completedConfig, cloud)
	}
}

func (c *ControllerWrapper) startNodeController(ctx context.Context, initContext app.ControllerInitContext, completedConfig *cloudcontrollerconfig.CompletedConfig, cloud cloudprovider.Interface) (controller.Interface, bool, error) {
	nodeController := NewNodeController()
	go nodeController.Run(ctx)
	return nil, true, nil
}
