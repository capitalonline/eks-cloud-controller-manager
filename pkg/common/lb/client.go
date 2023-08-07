package eks

import (
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/profile"
)

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

func NewCreateLbInstanceRequest() (request *CreateLbInstanceRequest) {
	request = &CreateLbInstanceRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionCreateLb)
	request.SetDomain(consts.ApiHost)
	return
}

func NewCreateLbInstanceResponse() (response *CreateLbInstanceResponse) {
	response = &CreateLbInstanceResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func (c *Client) CreateLbInstance(request *CreateLbInstanceRequest) (response *CreateLbInstanceResponse, err error) {
	if request == nil {
		request = NewCreateLbInstanceRequest()
	}
	response = NewCreateLbInstanceResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeLbInstanceRequest() (request *DescribeLbInstanceRequest) {
	request = &DescribeLbInstanceRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbInstance)
	request.SetDomain(consts.ApiHost)
	return
}

func NewDescribeLbInstanceResponse() (response *DescribeLbInstanceResponse) {
	response = &DescribeLbInstanceResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DescribeLbInstance(request *DescribeLbInstanceRequest) (response *DescribeLbInstanceResponse, err error) {
	if request == nil {
		request = NewDescribeLbInstanceRequest()
	}
	response = NewDescribeLbInstanceResponse()
	err = c.Send(request, response)
	return
}

func NewUpdateLbInstanceRequest() (request *UpdateLbInstanceRequest) {
	request = &UpdateLbInstanceRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionUpdateLbInstance)
	request.SetDomain(consts.ApiHost)
	return
}

func NewUpdateLbInstanceResponse() (response *UpdateLbInstanceResponse) {
	response = &UpdateLbInstanceResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) UpdateLbInstance(request *UpdateLbInstanceRequest) (response *UpdateLbInstanceResponse, err error) {
	if request == nil {
		request = NewUpdateLbInstanceRequest()
	}
	response = NewUpdateLbInstanceResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeTaskStatusRequest() (request *DescribeTaskStatusRequest) {
	request = &DescribeTaskStatusRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbTaskStatus)
	request.SetDomain(consts.ApiHost)
	return
}

func NewDescribeTaskStatusResponse() (response *DescribeTaskStatusResponse) {
	response = &DescribeTaskStatusResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DescribeTaskStatus(request *DescribeTaskStatusRequest) (response *DescribeTaskStatusResponse, err error) {
	if request == nil {
		request = NewDescribeTaskStatusRequest()
	}
	response = NewDescribeTaskStatusResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteLbListenersRequest() (request *DeleteLbListenersRequest) {
	request = &DeleteLbListenersRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDeleteLbInstance)
	request.SetDomain(consts.ApiHost)
	return
}

func NewDeleteLbListenersResponse() (response *DeleteLbListenersResponse) {
	response = &DeleteLbListenersResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DeleteLbListeners(request *DeleteLbListenersRequest) (response *DeleteLbListenersResponse, err error) {
	if request == nil {
		request = NewDeleteLbListenersRequest()
	}
	response = NewDeleteLbListenersResponse()
	err = c.Send(request, response)
	return
}
