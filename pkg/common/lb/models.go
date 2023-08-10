package lb

import (
	"encoding/json"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
)

type PackageCreateSlbRequest struct {
	*cdshttp.BaseRequest
	UserId            string                        `json:"UserId"`
	CustomerId        string                        `json:"CustomerId"`
	VpcId             string                        `json:"VpcId"`
	AvailableZoneCode string                        `json:"AvailableZoneCode"`
	Level             int                           `json:"Level"`
	SlbInfo           PackageCreateSlbInfo          `json:"SlbInfo"`
	BandwidthInfo     PackageCreateSlbBandwidthInfo `json:"BandwidthInfo"`
}

type PackageCreateSlbInfo struct {
	BillingSchemeId string `json:"BillingSchemeId"`
	NetType         string `json:"NetType"`
	Name            string `json:"Name"`
	SubjectId       int    `json:"SubjectId"`
}

type PackageCreateSlbBandwidthInfo struct {
	Name            string `json:"Name"`
	BillingSchemeId string `json:"BillingSchemeId"`
	Qos             int    `json:"Qos"`
	Type            string `json:"Type"`
	SubjectId       int    `json:"SubjectId"`
	IsAutoRenewal   bool   `json:"IsAutoRenewal"`
	IsToMonth       bool   `json:"IsToMonth"`
	Duration        int    `json:"Duration"`
	EipCount        int    `json:"EipCount"`
}

func (r *PackageCreateSlbRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *PackageCreateSlbRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type PackageCreateSlbResponse struct {
	*cdshttp.BaseResponse
	Data      []PackageCreateSlbResponseData `json:"Data"`
	Code      string                         `json:"Code"`
	Message   string                         `json:"Message"`
	Success   interface{}                    `json:"Success"`
	RequestId string                         `json:"RequestId"`
}

type PackageCreateSlbResponseData struct {
	SlbId  string `json:"SlbId"`
	TaskId string `json:"TaskId"`
}

func (r *PackageCreateSlbResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *PackageCreateSlbResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeVpcSlbRequest struct {
	*cdshttp.BaseRequest
	SlbID   string `json:"SlbId,omitempty"`
	SlbName string `json:"SlbName,omitempty"`
}

func (r *DescribeVpcSlbRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeVpcSlbRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeVpcSlbResponse struct {
	*cdshttp.BaseResponse

	Data      []DescribeVpcSlbResponseData `json:"Data"`
	Code      string                       `json:"Code"`
	Message   string                       `json:"Message"`
	Success   interface{}                  `json:"Success"`
	RequestId string                       `json:"RequestId"`
}

type DescribeVpcSlbResponseData struct {
	SlbInfo DescribeVpcSlbResponseSlbInfo `json:"SlbInfo"`
}

type DescribeVpcSlbResponseSlbInfo struct {
	SlbId         string                          `json:"SlbId"`
	SlbName       string                          `json:"SlbName"`
	SlbStatus     string                          `json:"SlbStatus"`
	BandwidthId   string                          `json:"BandwidthId"`
	BandwidthName string                          `json:"BandwidthName"`
	Qos           int                             `json:"Qos"`
	VipList       []DescribeVpcSlbResponseVipInfo `json:"VipList"`
}

type DescribeVpcSlbResponseVipInfo struct {
	VipId      string                           `json:"VipId"`
	VipIp      string                           `json:"VipIp"`
	VipStatus  string                           `json:"VipStatus"`
	ListenInfo DescribeVpcSlbResponseListenInfo `json:"ListenInfo"`
}

type DescribeVpcSlbResponseListenInfo struct {
	ListenId       string                         `json:"ListenId"`
	ListenPort     string                         `json:"ListenPort"`
	ListenProtocol string                         `json:"ListenProtocol"`
	RsList         []DescribeVpcSlbResponseRsInfo `json:"RsList"`
}

type DescribeVpcSlbResponseRsInfo struct {
	RsIp   string `json:"RsIp"`
	RsPort string `json:"RsPort"`
}

func (r *DescribeVpcSlbResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeVpcSlbResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbUpdateListenRequest struct {
	*cdshttp.BaseRequest
	SlbId      string                            `json:"SlbId"`
	Platform   string                            `json:"Platform"`
	ListenList []VpcSlbUpdateListenRequestListen `json:"ListenList"`
}

type VpcSlbUpdateListenRequestListen struct {
	ListenIp       string                        `json:"ListenIp"`
	ListenPort     int                           `json:"ListenPort"`
	ListenProtocol string                        `json:"ListenProtocol"`
	RsList         []VpcSlbUpdateListenRequestRs `json:"RsList"`
}

type VpcSlbUpdateListenRequestRs struct {
	RsIp   string `json:"RsIp"`
	RsPort string `json:"RsPort"`
}

func (r *VpcSlbUpdateListenRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbUpdateListenRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbUpdateListenResponse struct {
	*cdshttp.BaseResponse
	Data      []VpcSlbUpdateListenResponseData `json:"Data"`
	Code      string                           `json:"Code"`
	Message   string                           `json:"Message"`
	Success   interface{}                      `json:"Success"`
	RequestId string                           `json:"RequestId"`
}

type VpcSlbUpdateListenResponseData struct {
	TaskId string `json:"TaskId"`
}

func (r *VpcSlbUpdateListenResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbUpdateListenResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbClearListenRequest struct {
	*cdshttp.BaseRequest
	SlbId string `json:"SlbId"`
}

func (r *VpcSlbClearListenRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbClearListenRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbClearListenResponse struct {
	*cdshttp.BaseResponse
	Data      []VpcSlbClearListenResponseData `json:"Data"`
	Code      string                          `json:"Code"`
	Message   string                          `json:"Message"`
	Success   interface{}                     `json:"Success"`
	RequestId string                          `json:"RequestId"`
}

type VpcSlbClearListenResponseData struct {
	TaskId string `json:"TaskId"`
}

func (r *VpcSlbClearListenResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbClearListenResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskRequest struct {
	*cdshttp.BaseRequest
	TaskId string `json:"TaskId"`
}

func (r *DescribeTaskRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeTaskRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskResponse struct {
	*cdshttp.BaseResponse
	Data      []DescribeTaskResponseData `json:"Data"`
	Code      string                     `json:"Code"`
	Message   string                     `json:"Message"`
	Success   interface{}                `json:"Success"`
	RequestId string                     `json:"RequestId"`
}

type DescribeTaskResponseData struct {
	TaskStatus string `json:"TaskStatus"`
}

func (r *DescribeTaskResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeTaskResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbBillingSchemeRequest struct {
	*cdshttp.BaseRequest
	AvailableZoneCode string `json:"AvailableZoneCode"`
	BillingMethod     string `json:"BillingMethod"`
	NetType           string `json:"NetType"`
}

func (r *VpcSlbBillingSchemeRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbBillingSchemeRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcSlbBillingSchemeResponse struct {
	*cdshttp.BaseResponse
	Data      []VpcSlbBillingSchemeResponseData `json:"Data"`
	Code      string                            `json:"Code"`
	Message   string                            `json:"Message"`
	Success   interface{}                       `json:"Success"`
	RequestId string                            `json:"RequestId"`
}

type VpcSlbBillingSchemeResponseData struct {
	GoodsId           int                                      `json:"GoodsId"`
	BillingMethod     string                                   `json:"BillingMethod"`
	BillingType       string                                   `json:"BillingType"`
	BillingSchemeId   string                                   `json:"BillingSchemeId"`
	BillingSchemeName string                                   `json:"BillingSchemeName"`
	BillingCycleId    string                                   `json:"BillingCycleId"`
	BillingItems      []VpcSlbBillingSchemeResponseBillingItem `json:"BillingItems"`
	ConfId            int                                      `json:"ConfId"`
	ConfName          string                                   `json:"ConfName"`
	Description       string                                   `json:"Description"`
}

type VpcSlbBillingSchemeResponseBillingItem struct {
	Key    string `json:"Key"`
	Id     string `json:"Id"`
	Name   string `json:"Name"`
	AttrId string `json:"AttrId"`
	Size   int    `json:"Size"`
}

func (r *VpcSlbBillingSchemeResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *VpcSlbBillingSchemeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}
