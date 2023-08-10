package lb

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

func NewPackageCreateSlbRequest() (request *PackageCreateSlbRequest) {
	request = &PackageCreateSlbRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionPackageCreateSlb)
	request.SetDomain(consts.ApiHost)
	return
}

func NewPackageCreateSlbResponse() (response *PackageCreateSlbResponse) {
	response = &PackageCreateSlbResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func (c *Client) PackageCreateSlb(request *PackageCreateSlbRequest) (response *PackageCreateSlbResponse, err error) {
	if request == nil {
		request = NewPackageCreateSlbRequest()
	}
	response = NewPackageCreateSlbResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeVpcSlbRequest() (request *DescribeVpcSlbRequest) {
	request = &DescribeVpcSlbRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbInstance)
	request.SetDomain(consts.ApiHost)
	return
}

func NewDescribeVpcSlbResponse() (response *DescribeVpcSlbResponse) {
	response = &DescribeVpcSlbResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcSlb(request *DescribeVpcSlbRequest) (response *DescribeVpcSlbResponse, err error) {
	if request == nil {
		request = NewDescribeVpcSlbRequest()
	}
	response = NewDescribeVpcSlbResponse()
	err = c.Send(request, response)
	return
}

func NewVpcSlbUpdateListenRequest() (request *VpcSlbUpdateListenRequest) {
	request = &VpcSlbUpdateListenRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionVpcSlbUpdateListen)
	request.SetDomain(consts.ApiHost)
	return
}

func NewVpcSlbUpdateListenResponse() (response *VpcSlbUpdateListenResponse) {
	response = &VpcSlbUpdateListenResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) VpcSlbUpdateListen(request *VpcSlbUpdateListenRequest) (response *VpcSlbUpdateListenResponse, err error) {
	if request == nil {
		request = NewVpcSlbUpdateListenRequest()
	}
	response = NewVpcSlbUpdateListenResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeTaskRequest() (request *DescribeTaskRequest) {
	request = &DescribeTaskRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbTaskStatus)
	request.SetDomain(consts.ApiHost)
	return
}

func NewDescribeTaskResponse() (response *DescribeTaskResponse) {
	response = &DescribeTaskResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DescribeTask(request *DescribeTaskRequest) (response *DescribeTaskResponse, err error) {
	if request == nil {
		request = NewDescribeTaskRequest()
	}
	response = NewDescribeTaskResponse()
	err = c.Send(request, response)
	return
}

//func NewDeleteLbListenersRequest() (request *DeleteLbListenersRequest) {
//	request = &DeleteLbListenersRequest{
//		BaseRequest: &cdshttp.BaseRequest{},
//	}
//	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDeleteLbInstance)
//	request.SetDomain(consts.ApiHost)
//	return
//}
//
//func NewDeleteLbListenersResponse() (response *DeleteLbListenersResponse) {
//	response = &DeleteLbListenersResponse{
//		BaseResponse: &cdshttp.BaseResponse{},
//	}
//	return
//}
//
//func (c *Client) DeleteLbListeners(request *DeleteLbListenersRequest) (response *DeleteLbListenersResponse, err error) {
//	if request == nil {
//		request = NewDeleteLbListenersRequest()
//	}
//	response = NewDeleteLbListenersResponse()
//	err = c.Send(request, response)
//	return
//}

func NewVpcSlbBillingSchemeRequest() (request *VpcSlbBillingSchemeRequest) {
	request = &VpcSlbBillingSchemeRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbTaskStatus)
	request.SetDomain(consts.ApiHost)
	return
}

func NewVpcSlbBillingSchemeResponse() (response *VpcSlbBillingSchemeResponse) {
	response = &VpcSlbBillingSchemeResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) VpcSlbBillingScheme(request *VpcSlbBillingSchemeRequest) (response *VpcSlbBillingSchemeResponse, err error) {
	if request == nil {
		request = NewVpcSlbBillingSchemeRequest()
	}
	response = NewVpcSlbBillingSchemeResponse()
	err = c.Send(request, response)
	return
}

func NewVpcSlbClearListenRequest() (request *VpcSlbClearListenRequest) {
	request = &VpcSlbClearListenRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionVpcSlbClearListen)
	request.SetDomain(consts.ApiHost)
	return
}

func NewVpcSlbClearListenResponse() (response *VpcSlbClearListenResponse) {
	response = &VpcSlbClearListenResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) VpcSlbClearListen(request *VpcSlbClearListenRequest) (response *VpcSlbClearListenResponse, err error) {
	if request == nil {
		request = NewVpcSlbClearListenRequest()
	}
	response = NewVpcSlbClearListenResponse()
	err = c.Send(request, response)
	return
}
