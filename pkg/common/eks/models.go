package eks

import (
	"encoding/json"
	cdshttp "github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/http"
)

type DescribeEKSNodeRequest struct {
	*cdshttp.BaseRequest
	ClusterId string `json:"ClusterId,omitempty"`
	NodeName  string `json:"NodeName,omitempty"`
	NodeId    string `json:"NodeId,omitempty"`
}

func (req *DescribeEKSNodeRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *DescribeEKSNodeRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type DescribeEKSNodeResponse struct {
	*cdshttp.BaseResponse
	Code string                       `json:"Code"`
	Msg  string                       `json:"Msg"`
	Data *DescribeEKSNodeResponseData `json:"Data"`
}

type DescribeEKSNodeResponseData struct {
	NodeId string                             `json:"NodeId"`
	Labels []DescribeEKSNodeResponseDataLabel `json:"Labels"`
	Taints []DescribeEKSNodeResponseDataTaint `json:"Taints"`
}

type DescribeEKSNodeResponseDataTaint struct {
	Key    string `json:"Key"`
	Value  string `json:"Value"`
	Effect string `json:"Effect"`
}

type DescribeEKSNodeResponseDataLabel struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (resp *DescribeEKSNodeResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *DescribeEKSNodeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

type NodeCCMInitRequest struct {
	*cdshttp.BaseRequest
	ClusterId string `json:"ClusterId,omitempty"`
	NodeName  string `json:"NodeName,omitempty"`
	NodeId    string `json:"NodeId,omitempty"`
}

func (req *NodeCCMInitRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *NodeCCMInitRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type NodeCCMInitResponse struct {
	*cdshttp.BaseResponse
	Code string                   `json:"Code"`
	Msg  string                   `json:"Msg"`
	Data *NodeCCMInitResponseData `json:"Data"`
}

type NodeCCMInitResponseData struct {
	PrivateIp   string       `json:"PrivateIp"`
	Status      string       `json:"Status"`
	Labels      []Label      `json:"Labels"`
	Taints      []Taint      `json:"Taints"`
	Annotations []Annotation `json:"Annotations"`
	NodeId      string       `json:"NodeId"`
}

type Label struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type Taint struct {
	Key    string `json:"Key"`
	Value  string `json:"Value"`
	Effect string `json:"Effect"`
}

type Annotation struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (resp *NodeCCMInitResponseData) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *NodeCCMInitResponseData) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

type ModifyClusterLoadRequest struct {
	*cdshttp.BaseRequest
	ClusterId string                     `json:"ClusterId,omitempty"`
	NodeList  []ModifyClusterLoadReqNode `json:"NodeList,omitempty"`
}

func (req *ModifyClusterLoadRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *ModifyClusterLoadRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type ModifyClusterLoadResponse struct {
	*cdshttp.BaseResponse
	Code string `json:"Code"`
	Msg  string `json:"Msg"`
}

type ModifyClusterLoadReqNode struct {
	NodeId   string        `json:"NodeId"`
	Cpu      *ResourceInfo `json:"Cpu,omitempty"`
	Memory   *ResourceInfo `json:"Memory,omitempty"`
	Status   string        `json:"Status"`
	NodeName string        `json:"-"`
}

type NodeLoad struct {
	Cpu    ResourceInfo `json:"Cpu"`
	Mem    ResourceInfo `json:"Mem"`
	Status string       `json:"Status"`
}

type ResourceInfo struct {
	Usage    int64 `json:"Usage"`
	Limits   int64 `json:"Limits"`
	Requests int64 `json:"Requests"`
}

func (resp *ModifyClusterLoadResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *ModifyClusterLoadResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}

type SendAlarmRequest struct {
	*cdshttp.BaseRequest
	Theme     string        `json:"Theme"`
	ClusterId string        `json:"ClusterId"`
	NodeId    string        `json:"NodeId"`
	Source    string        `json:"Source"`
	Keyword   string        `json:"Keyword"`
	Metric    string        `json:"Metric"`
	Value     interface{}   `json:"Value"`
	Tags      []interface{} `json:"Tags"`
	AlarmMsg  string        `json:"AlarmMsg"`
}

func (req *SendAlarmRequest) ToJsonString() string {
	b, _ := json.Marshal(req)
	return string(b)
}

func (req *SendAlarmRequest) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &req)
}

type SendAlarmResponse struct {
	*cdshttp.BaseResponse
	Code string `json:"Code"`
	Msg  string `json:"Msg"`
}

func (resp *SendAlarmResponse) ToJsonString() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func (resp *SendAlarmResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &resp)
}
