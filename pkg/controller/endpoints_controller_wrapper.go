package controller

import (
	"context"
	"fmt"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	genericcontrollermanager "k8s.io/controller-manager/app"
	"k8s.io/controller-manager/controller"
)

const (
	EndpointsControllerName = "endpoints-controller"
)

// EndpointsControllerWrapper 包装 EndpointsController，用于注册到 ControllerInitializers
type EndpointsControllerWrapper struct{}

func NewEndpointsControllerWrapper() *EndpointsControllerWrapper {
	return &EndpointsControllerWrapper{}
}

// StartEndpointsControllerWrapper 初始化并启动 EndpointsController
func (e *EndpointsControllerWrapper) StartEndpointsControllerWrapper(
	initContext app.ControllerInitContext,
	completedConfig *cloudcontrollerconfig.CompletedConfig,
	cloud cloudprovider.Interface,
) app.InitFunc {
	return func(ctx context.Context, controllerContext genericcontrollermanager.ControllerContext) (controller.Interface, bool, error) {
		endpointsController, err := NewServiceController()
		if err != nil {
			return nil, false, fmt.Errorf("failed to create EndpointsController: %v", err)
		}
		go endpointsController.Run(ctx)
		return nil, true, nil
	}
}
