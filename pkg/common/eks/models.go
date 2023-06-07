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
