package api

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/profile"
	"net/http"
)

func NodeAddresses(clusterId, nodeId string, nodeName string) ([]string, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodGet
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	request := eks.NewNodeCCMInitRequest()
	request.NodeId = nodeId
	request.ClusterId = clusterId
	request.NodeName = nodeName
	response, err := client.NodeCCMInit(request)
	if err != nil {
		return nil, err
	}
	return []string{response.Data.PrivateIp}, err
}

func NodeCCMInit(clusterId, nodeId string, nodeName string) (*eks.NodeCCMInitResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodGet
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	request := eks.NewNodeCCMInitRequest()
	request.NodeId = nodeId
	request.ClusterId = clusterId
	request.NodeName = nodeName
	response, err := client.NodeCCMInit(request)
	if err != nil {
		return nil, err
	}

	return response, err
}

func ModifyClusterLoad(request *eks.ModifyClusterLoadRequest) (*eks.ModifyClusterLoadResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	response, err := client.ModifyClusterLoad(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func NotifyMasterDown(request *eks.SendAlarmRequest) (*eks.SendAlarmResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	response, err := client.SendAlarm(request)
	if err != nil {
		return nil, err
	}
	return response, err
}
