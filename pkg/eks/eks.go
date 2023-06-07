package eks

import (
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/profile"
	"net/http"
)

func NodeAddresses(clusterId, nodeId, nodeName string) ([]string, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	request := eks.NewDescribeEKSNodeRequest()
	request.NodeId = nodeId
	request.ClusterId = clusterId
	request.NodeName = nodeName
	response, err := client.DescribeEKSNode(request)
	if err != nil {
		return nil, err
	}
	fmt.Println(response)
	return []string{}, err
}

func DescribeNodeDetails(clusterId, nodeId string) (*eks.DescribeEKSNodeResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	client, _ := eks.NewClient(credential, consts.Region, cpf)
	request := eks.NewDescribeEKSNodeRequest()
	request.NodeId = nodeId
	request.ClusterId = clusterId
	response, err := client.DescribeEKSNode(request)
	if err != nil {
		return nil, err
	}

	// TODO 去掉此处逻辑
	response = &eks.DescribeEKSNodeResponse{
		Code: "Success",
		Msg:  "",
		Data: &eks.DescribeEKSNodeResponseData{
			NodeId: nodeId,
			Labels: []eks.DescribeEKSNodeResponseDataLabel{
				{
					Key:   "",
					Value: "",
				},
			},
			Taints: []eks.DescribeEKSNodeResponseDataTaint{{
				Key:   "",
				Value: "",
			}},
		},
	}

	fmt.Println(response)
	return response, err
}

//func ()  {
//
//}
