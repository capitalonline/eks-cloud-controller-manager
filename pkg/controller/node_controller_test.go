package controller

import (
	"context"
	_ "github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"testing"
)

func TestNewNodeController(t *testing.T) {
	var nodeController = NewNodeController()
	//nodeController.CollectPlayLoad(context.Background())
	nodeController.Update(context.Background())
}
