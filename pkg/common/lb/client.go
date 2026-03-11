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
	request.SetDomain(consts.LbApiHost)
	return
}

func NewStandardCreateSlbRequest() (request *StandardCreateSlbRequest) {
	request = &StandardCreateSlbRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionStandardCreateVpcSlb)
	request.SetDomain(consts.LbApiHost)
	return
}

func NewPackageCreateSlbResponse() (response *PackageCreateSlbResponse) {
	response = &PackageCreateSlbResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

func NewStandardCreateSlbResponse() (response *StandardCreateSlbResponse) {
	response = &StandardCreateSlbResponse{BaseResponse: &cdshttp.BaseResponse{}}
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

func (c *Client) StandardCreateSlb(request *StandardCreateSlbRequest) (response *StandardCreateSlbResponse, err error) {
	if request == nil {
		request = NewStandardCreateSlbRequest()
	}
	response = NewStandardCreateSlbResponse()
	err = c.Send(request, response)
	return
}
func NewDescribeVpcSlbRequest() (request *DescribeVpcSlbRequest) {
	request = &DescribeVpcSlbRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDescribeLbInstance)
	request.SetDomain(consts.LbApiHost)
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
	request.SetDomain(consts.LbApiHost)
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
	request.SetDomain(consts.LbApiHost)
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

func NewDeleteLbListenersRequest() (request *DeleteVpcSLBListenRequest) {
	request = &DeleteVpcSLBListenRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDeleteVpcSLBListen)
	request.SetDomain(consts.LbApiHost)
	return
}

func NewDeleteLbListenersResponse() (response *DeleteVpcSLBListenResponse) {
	response = &DeleteVpcSLBListenResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DeleteVpcSLBListen(request *DeleteVpcSLBListenRequest) (response *DeleteVpcSLBListenResponse, err error) {
	if request == nil {
		request = NewDeleteLbListenersRequest()
	}
	response = NewDeleteLbListenersResponse()
	err = c.Send(request, response)
	return
}

func NewVpcSlbBillingSchemeRequest() (request *VpcSlbBillingSchemeRequest) {
	request = &VpcSlbBillingSchemeRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionVpcSlbBillingScheme)
	request.SetDomain(consts.LbApiHost)
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
	request.SetDomain(consts.LbApiHost)
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

func NewBandwidthBillingSchemeRequest() (request *BandwidthBillingSchemeRequest) {
	request = &BandwidthBillingSchemeRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionVpcBandwidthBillingScheme)
	request.SetDomain(consts.LbApiHost)
	return
}

func NewBandwidthBillingSchemeResponse() (response *BandwidthBillingSchemeResponse) {
	response = &BandwidthBillingSchemeResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) BandwidthBillingScheme(request *BandwidthBillingSchemeRequest) (response *BandwidthBillingSchemeResponse, err error) {
	if request == nil {
		request = NewBandwidthBillingSchemeRequest()
	}
	response = NewBandwidthBillingSchemeResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteVpcSlbRequest() (request *DeleteVpcSlbRequest) {
	request = &DeleteVpcSlbRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.Init().WithApiInfo(consts.ServiceLb, consts.ApiVersion, consts.ActionDeleteVpcSlb)
	request.SetDomain(consts.LbApiHost)
	return
}

func NewDeleteVpcSlbResponse() (response *DeleteVpcSlbResponse) {
	response = &DeleteVpcSlbResponse{
		BaseResponse: &cdshttp.BaseResponse{},
	}
	return
}

func (c *Client) DeleteVpcSlb(request *DeleteVpcSlbRequest) (response *DeleteVpcSlbResponse, err error) {
	if request == nil {
		request = NewDeleteVpcSlbRequest()
	}
	response = NewDeleteVpcSlbResponse()
	err = c.Send(request, response)
	return
}
