package api

import (
	"errors"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/lb"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/utils/profile"
	"net/http"
)

func DescribeVpcSlb(request *lb.DescribeVpcSlbRequest) (*lb.DescribeVpcSlbResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodGet
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := lb.NewClient(credential, consts.Region, cpf)
	response, err := client.DescribeVpcSlb(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func PackageCreateSlb(request *lb.PackageCreateSlbRequest) (*lb.PackageCreateSlbResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := lb.NewClient(credential, consts.Region, cpf)
	response, err := client.PackageCreateSlb(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func VpcSlbBillingScheme(request *lb.VpcSlbBillingSchemeRequest) (*lb.VpcSlbBillingSchemeResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodGet
	cpf.HttpProfile.Endpoint = consts.ApiHost
	client, _ := lb.NewClient(credential, consts.Region, cpf)
	response, err := client.VpcSlbBillingScheme(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func DescribeTask(taskId string) (*lb.DescribeTaskResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodGet
	cpf.HttpProfile.Endpoint = consts.ApiHost
	request := lb.NewDescribeTaskRequest()
	request.TaskId = taskId
	client, _ := lb.NewClient(credential, consts.Region, cpf)
	response, err := client.DescribeTask(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func VpcSlbClearListen(slbId string) (*lb.VpcSlbClearListenResponse, error) {
	credential := utils.NewCredential(consts.AccessKeyID, consts.AccessKeySecret)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	cpf.HttpProfile.Endpoint = consts.ApiHost
	request := lb.NewVpcSlbClearListenRequest()
	request.SlbId = slbId
	client, _ := lb.NewClient(credential, consts.Region, cpf)
	response, err := client.VpcSlbClearListen(request)
	if err != nil {
		return nil, err
	}
	if response == nil || len(response.Data) < 1 {
		return nil, errors.New("清空监听规则接口错误")
	}
	return response, err
}
