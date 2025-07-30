package eks

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/profile"
)

// test commit

type Client struct {
	utils.Client
}

func NewClient(credential *utils.Credential, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
	client = &Client{}
	client.Init(region).
		WithCredential(credential).
		WithProfile(clientProfile)
	return
}

func NewDescribeEKSNodeRequest() (request *DescribeEKSNodeRequest) {
	request = &DescribeEKSNodeRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionDescribeEKSNode)
	return
}

func NewDescribeEKSNodeResponse() (response *DescribeEKSNodeResponse) {
	response = &DescribeEKSNodeResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func NewNodeCCMInitRequest() (request *NodeCCMInitRequest) {
	request = &NodeCCMInitRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionNodeCCMInit)
	return
}

func NewNodeCCMInitResponse() (response *NodeCCMInitResponse) {
	response = &NodeCCMInitResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func NewModifyClusterLoadRequest() (request *ModifyClusterLoadRequest) {
	request = &ModifyClusterLoadRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionModifyClusterLoad)
	return
}

func NewModifyClusterLoadResponse() (response *ModifyClusterLoadResponse) {
	response = &ModifyClusterLoadResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func (c *Client) DescribeEKSNode(request *DescribeEKSNodeRequest) (response *DescribeEKSNodeResponse, err error) {
	if request == nil {
		request = NewDescribeEKSNodeRequest()
	}
	response = NewDescribeEKSNodeResponse()
	err = c.Send(request, response)
	return
}

func (c *Client) NodeCCMInit(request *NodeCCMInitRequest) (response *NodeCCMInitResponse, err error) {
	if request == nil {
		request = NewNodeCCMInitRequest()
	}
	response = NewNodeCCMInitResponse()
	err = c.Send(request, response)
	return
}

func (c *Client) ModifyClusterLoad(request *ModifyClusterLoadRequest) (response *ModifyClusterLoadResponse, err error) {
	if request == nil {
		request = NewModifyClusterLoadRequest()
	}
	response = NewModifyClusterLoadResponse()
	err = c.Send(request, response)
	return
}

func (c *Client) SendAlarm(request *SendAlarmRequest) (response *SendAlarmResponse, err error) {
	if request == nil {
		request = NewSendAlarmRequest()
	}
	response = NewSendAlarmResponse()
	err = c.Send(request, response)
	return
}

func NewSendAlarmRequest() (request *SendAlarmRequest) {
	request = &SendAlarmRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionSendAlarm)
	return
}

func NewSendAlarmResponse() (response *SendAlarmResponse) {
	response = &SendAlarmResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}
