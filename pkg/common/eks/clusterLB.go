package eks

import (
	"encoding/json"

	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
)

// RegisterClusterLB 注册EKS集群LB信息
func (c *Client) RegisterClusterLB(request *RegisterClusterLBRequest) (response *RegisterClusterLBResponse, err error) {
	if request == nil {
		request = NewRegisterClusterLBRequest()
	}
	response = NewRegisterClusterLBResponse()
	err = c.Send(request, response)
	return
}

func NewRegisterClusterLBRequest() (request *RegisterClusterLBRequest) {
	request = &RegisterClusterLBRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionCreateLBRegistration)
	return
}

func NewRegisterClusterLBResponse() (response *RegisterClusterLBResponse) {
	response = &RegisterClusterLBResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

type RegisterClusterLBRequest struct {
	*cdshttp.BaseRequest
	ClusterId string            `json:"ClusterId"`
	SvcName   string            `json:"SvcName"`
	SvcNs     string            `json:"SvcNs"`
	SlbId     string            `json:"SlbId"`
	Pattern   string            `json:"Pattern"`
	SvcUid    string            `json:"SvcUid"`
	Eip       string            `json:"Eip"`
	Vip       string            `json:"Vip"`
	Ports     string            `json:"Ports"`
	Info      map[string]string `json:"Info"`
}

func (req *RegisterClusterLBRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *RegisterClusterLBRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type RegisterClusterLBResponse struct {
	*cdshttp.BaseResponse
	Code string `json:"Code"`
	Msg  string `json:"Msg"`
}

func (resp *RegisterClusterLBResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *RegisterClusterLBResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

// UpdateClusterLB 更新EKS集群LB信息
func (c *Client) UpdateClusterLB(request *UpdateClusterLBRequest) (response *UpdateClusterLBResponse, err error) {
	if request == nil {
		request = NewUpdateClusterLBRequest()
	}
	response = NewUpdateClusterLBResponse()
	err = c.Send(request, response)
	return
}

func NewUpdateClusterLBRequest() (request *UpdateClusterLBRequest) {
	request = &UpdateClusterLBRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionUpdateLBPorts)
	return
}

func NewUpdateClusterLBResponse() (response *UpdateClusterLBResponse) {
	response = &UpdateClusterLBResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

type UpdateClusterLBRequest struct {
	*cdshttp.BaseRequest
	CustomerId string `json:"CustomerId"`
	UserId     string `json:"UserId"`
	ClusterId  string `json:"ClusterId"`
	SvcUid     string `json:"SvcUid"`
	SvcName    string `json:"SvcName"`
	SvcNs      string `json:"SvcNs"`
	SlbId      string `json:"SlbId"`
	Ports      string `json:"Ports"`
}

func (req *UpdateClusterLBRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *UpdateClusterLBRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type UpdateClusterLBResponse struct {
	*cdshttp.BaseResponse
	Code string `json:"Code"`
	Msg  string `json:"Msg"`
}

func (resp *UpdateClusterLBResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *UpdateClusterLBResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

// DeleteClusterLB 删除EKS集群LB信息
func (c *Client) DeleteClusterLB(request *DeleteClusterLBRequest) (response *DeleteClusterLBResponse, err error) {
	if request == nil {
		request = NewDeleteClusterLBRequest()
	}
	response = NewDeleteClusterLBResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteClusterLBRequest() (request *DeleteClusterLBRequest) {
	request = &DeleteClusterLBRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionDeleteLBRegistration)
	return
}

func NewDeleteClusterLBResponse() (response *DeleteClusterLBResponse) {
	response = &DeleteClusterLBResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

type DeleteClusterLBRequest struct {
	*cdshttp.BaseRequest
	CustomerId string `json:"CustomerId"`
	UserId     string `json:"UserId"`
	ClusterId  string `json:"ClusterId"`
	SvcUid     string `json:"SvcUid"`
	SvcName    string `json:"SvcName"`
	SvcNs      string `json:"SvcNs"`
}

func (req *DeleteClusterLBRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *DeleteClusterLBRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type DeleteClusterLBResponse struct {
	*cdshttp.BaseResponse
	Code string `json:"Code"`
	Msg  string `json:"Msg"`
}

func (resp *DeleteClusterLBResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *DeleteClusterLBResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

// GetClusterLBDisabledPorts 获取EKS集群LB失效端口
func (c *Client) GetClusterLBDisabledPorts(request *GetClusterLBDisabledPortsRequest) (response *GetClusterLBDisabledPortsResponse, err error) {
	if request == nil {
		request = NewGetClusterLBDisabledPortsRequest()
	}
	response = NewGetClusterLBDisabledPortsResponse()
	err = c.Send(request, response)
	return
}

func NewGetClusterLBDisabledPortsRequest() (request *GetClusterLBDisabledPortsRequest) {
	request = &GetClusterLBDisabledPortsRequest{
		BaseRequest: &cdshttp.BaseRequest{},
	}
	request.SetDomain(consts.ApiHost)
	request.Init().WithApiInfo(consts.ServiceEKS, consts.ApiVersion, consts.ActionQueryInvalidLBPorts)
	return
}

func NewGetClusterLBDisabledPortsResponse() (response *GetClusterLBDisabledPortsResponse) {
	response = &GetClusterLBDisabledPortsResponse{BaseResponse: &cdshttp.BaseResponse{}}
	return
}

type GetClusterLBDisabledPortsRequest struct {
	*cdshttp.BaseRequest
	CustomerId string `json:"CustomerId"`
	UserId     string `json:"UserId"`
	ClusterId  string `json:"ClusterId"`
	SvcUid     string `json:"SvcUid"`
	SvcName    string `json:"SvcName"`
	SvcNs      string `json:"SvcNs"`
	SlbId      string `json:"SlbId"`
	Ports      string `json:"Ports"`
}

func (req *GetClusterLBDisabledPortsRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *GetClusterLBDisabledPortsRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type GetClusterLBDisabledPortsResponse struct {
	*cdshttp.BaseResponse
	Code string                 `json:"Code"`
	Msg  string                 `json:"Msg"`
	Data QueryInvalidLBPortsRes `json:"Data"`
}

type QueryInvalidLBPortsRes struct {
	InvalidPorts []int `json:"InvalidPorts"`
}

func (resp *GetClusterLBDisabledPortsResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *GetClusterLBDisabledPortsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}
