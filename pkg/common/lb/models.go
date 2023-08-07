package eks

import (
	"encoding/json"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
)

type CreateLbInstanceRequest struct {
	*cdshttp.BaseRequest
	UserId            string                           `json:"UserId"`
	CustomerId        string                           `json:"CustomerId"`
	VpcId             string                           `json:"VpcId"`
	AvailableZoneCode string                           `json:"AvailableZoneCode"`
	Level             int                              `json:"Level"`
	SlbInfo           CreateLbInstanceSlbInfo          `json:"SlbInfo"`
	BandwidthInfo     CreateLbInstanceSlbBandwidthInfo `json:"BandwidthInfo"`
}

type CreateLbInstanceSlbInfo struct {
	BillingSchemeId string `json:"BillingSchemeId"`
	NetType         string `json:"NetType"`
	Name            string `json:"Name"`
	SubjectId       int    `json:"SubjectId"`
}

type CreateLbInstanceSlbBandwidthInfo struct {
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

func (r *CreateLbInstanceRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *CreateLbInstanceRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CreateLbInstanceResponse struct {
	*cdshttp.BaseResponse
	Data      []CreateLbInstanceResponseData `json:"Data"`
	Code      string                         `json:"Code"`
	Message   string                         `json:"Message"`
	Success   interface{}                    `json:"Success"`
	RequestId string                         `json:"RequestId"`
}

type CreateLbInstanceResponseData struct {
	SlbId  string `json:"SlbId"`
	TaskId string `json:"TaskId"`
}

func (r *CreateLbInstanceResponseData) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *CreateLbInstanceResponseData) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeLbInstanceRequest struct {
	*cdshttp.BaseRequest
	SlbID string `json:"SlbId"`
}

func (r *DescribeLbInstanceRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeLbInstanceRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeLbInstanceResponse struct {
	*cdshttp.BaseResponse

	Data      []DescribeLbInstanceResponseData `json:"Data"`
	Code      string                           `json:"Code"`
	Message   string                           `json:"Message"`
	Success   interface{}                      `json:"Success"`
	RequestId string                           `json:"RequestId"`
}

type DescribeLbInstanceResponseData struct {
	SlbInfo DescribeLbInstanceResponseSlbInfo `json:"SlbInfo"`
}

type DescribeLbInstanceResponseSlbInfo struct {
	SlbId         string                              `json:"SlbId"`
	SlbName       string                              `json:"SlbName"`
	SlbStatus     string                              `json:"SlbStatus"`
	BandwidthId   string                              `json:"BandwidthId"`
	BandwidthName string                              `json:"BandwidthName"`
	Qos           int                                 `json:"Qos"`
	EipList       []DescribeLbInstanceResponseEipInfo `json:"EipList"`
}

type DescribeLbInstanceResponseEipInfo struct {
	EipId      string                               `json:"EipId"`
	EipIp      string                               `json:"EipIp"`
	EipStatus  string                               `json:"EipStatus"`
	ListenInfo DescribeLbInstanceResponseListenInfo `json:"ListenInfo"`
}

type DescribeLbInstanceResponseListenInfo struct {
	ListenId       string                             `json:"ListenId"`
	ListenPort     string                             `json:"ListenPort"`
	ListenProtocol string                             `json:"ListenProtocol"`
	RsList         []DescribeLbInstanceResponseRsInfo `json:"RsList"`
}

type DescribeLbInstanceResponseRsInfo struct {
	RsIp   string `json:"RsIp"`
	RsPort string `json:"RsPort"`
}

func (r *DescribeLbInstanceResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeLbInstanceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type UpdateLbInstanceRequest struct {
	*cdshttp.BaseRequest
	SlbId      string                          `json:"SlbId"`
	Platform   string                          `json:"Platform"`
	ListenList []UpdateLbInstanceRequestListen `json:"ListenList"`
}

type UpdateLbInstanceRequestListen struct {
	ListenIp       string                      `json:"ListenIp"`
	ListenPort     int                         `json:"ListenPort"`
	ListenProtocol string                      `json:"ListenProtocol"`
	RsList         []UpdateLbInstanceRequestRs `json:"RsList"`
}

type UpdateLbInstanceRequestRs struct {
	RsIp   string `json:"RsIp"`
	RsPort string `json:"RsPort"`
}

func (r *UpdateLbInstanceRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *UpdateLbInstanceRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type UpdateLbInstanceResponse struct {
	*cdshttp.BaseResponse
	Data      []UpdateLbInstanceResponseData `json:"Data"`
	Code      string                         `json:"Code"`
	Message   string                         `json:"Message"`
	Success   interface{}                    `json:"Success"`
	RequestId string                         `json:"RequestId"`
}

type UpdateLbInstanceResponseData struct {
	TaskId string `json:"TaskId"`
}

func (r *UpdateLbInstanceResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *UpdateLbInstanceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeleteLbListenersRequest struct {
	*cdshttp.BaseRequest
	SlbId string `json:"SlbId"`
}

func (r *DeleteLbListenersRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DeleteLbListenersRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeleteLbListenersResponse struct {
	*cdshttp.BaseResponse
	Data      []DeleteLbListenersResponseData `json:"Data"`
	Code      string                          `json:"Code"`
	Message   string                          `json:"Message"`
	Success   interface{}                     `json:"Success"`
	RequestId string                          `json:"RequestId"`
}

type DeleteLbListenersResponseData struct {
	TaskId string `json:"TaskId"`
}

func (r *DeleteLbListenersResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DeleteLbListenersResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusRequest struct {
	*cdshttp.BaseRequest
	TaskId string `json:"TaskId"`
}

func (r *DescribeTaskStatusRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeTaskStatusRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusResponse struct {
	*cdshttp.BaseResponse
	Data      []DescribeTaskStatusResponseData `json:"Data"`
	Code      string                           `json:"Code"`
	Message   string                           `json:"Message"`
	Success   interface{}                      `json:"Success"`
	RequestId string                           `json:"RequestId"`
}

type DescribeTaskStatusResponseData struct {
	TaskStatus string `json:"TaskStatus"`
}

func (r *DescribeTaskStatusResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *DescribeTaskStatusResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}
