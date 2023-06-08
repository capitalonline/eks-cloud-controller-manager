package controller

import (
	"context"
	"errors"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	commoneks "github.com/capitalonline/eks-cloud-controller-manager/pkg/common/eks"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/eks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"log"
	"time"
)

type NodeController struct {
	clientSet *kubernetes.Clientset
}

func (n *NodeController) Validate() error {
	return nil
}

var s = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1EWXdPREF6TWpZMU5sb1hEVE16TURZd05UQXpNalkxTmxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTkpUCjBWTk4weWNIVHFnYlhaSUkzQVYyb3VEaFA5ckJQTXRUZVl6eHI0SDNJYWdSM0crNjg5bnM2dGQ1RFJSUUFCUmcKL0xHNlRhRTNIUzdWR1lZcVV3anlPbzVIZDlkdGhVZnlzUE11Vkx0M1VBL3J6aks5UXhPUTZ0VHEyRzZFcklOawpYa0t5OVVmSDVOTll2aEM2KzF0QXZGaTN6SGUyOTdKbE1qOEF2b25DblFLT20wa3ZFNFVaSnJPY2RPdUpVR1BhCkxrS2M2djExS3lRckp1ODdmNGU4SGs1T0NMS0wrUHZydjN3VFVLdllrTHp4cmRKbUVDNFBJaytmODd3ZFgrMW0KZkNGOGljcUpCbHFTUytUZkZHV3BsV3d4VlU1WUV1Z3p0bXJFbkxScDBKU21PUnlaSENZR2p1cFNYd3RzaGpLKwpXNWxLQ2t4MTZOaUt2b0pvRS8wQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZHamxNNmt2dXVwNzVHSVF5UFZiMm5GU2o2dFFNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBRitxY2pCWlRlQUVESGQ4OW9aVAp5d3M2Q2FYVW04Z3FZeVpaVmFWbTR1cnhsOE1Ebi9UNVYxR0loZmd6U3ZxVmhLYXl1RjB0ckNWU3B5RkVRMzRYCjBaQVZZWmR0Y01oS0xKTGJ3b2taRjc3Z1puSEtpTE1SeVJNYWRpTXZycnU4OUZXNTF3STZUdW54dFp6SEFmSGIKOVdtVlNLVTVXZVczQ2J0cXZsUHR2TEt5TWFzWlVlZlU4SVMrdElsd0laemtlb3YwRWI0dWdJNWtQc05ZTkNlegpQYTZ2Slkvc0hFcDNjR0EzTWlaN1dSTjlnNkQwcUVxRzhxYlR0ZWFwYTQwUFFDaExDZm9JMWg3TmkvYVpXc0hxCmZGMnd3c2F1bGpqWUJpbEdJcERpaytoL1FUU2taakJQQ2ttY05HajIvWjlhR25GRExTeHliMGM3VUQrcjU2RVcKbmNvPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://192.168.137.183:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJSlh5bkltOElNZmd3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpBMk1EZ3dNekkyTlRaYUZ3MHlOREEyTURjd016STJOVGRhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXlSRmwvMWxCOE9RSmdkWDUKZzlYdzFYbmNmUERrOUN3WElFcTV6aFJ2KytQNHlDTDdEMTlvTmhyeHZmQzRWeS9PVEoreFoxc2k4TU0zcG1xagpCNERlV2FLWTZBQlVTcmxzVTQxcnUrQndCV3d3QWFLdHNOeG0vY0t5ZUV0UFcrRlhONkZEMWUvRnZIaVR4VXlhCnVYdTdpYmc1UHBoSVkwK0UxRTh4cFlWUTZkMzF3REtlaUQwRnc2UlBHVFh0aHliaC96Z1VHL2FVclNuZHI3UmUKMGRvLzI5ZHh3UzhkWGR6U3pWSnlkY08xTFQzUHhpK2ZmeXFOOERsSXlnU1kwdU1nUDFVaS96eXpDSDkvaHY0dgozeXNRZFBTeUxaMk5ubFBoY1NuTFljTlZLNk9HNEVMcW5LYTVBVlN1Q2VpL3BhVTVlTlptcFkyNXJWYlkxMkVvCldYNEo5UUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JSbzVUT3BMN3JxZStSaUVNajFXOXB4VW8rcgpVREFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBV2ZPeVlFYWtWLzloUU1wNVl3OXZFZXpEUk9ndEpwMmFvcy9VCmJjaFR4b21SYzRtMzZDSDZLNWhvS2JyU05tRVhLNEFIWVVBRzZZNkNKdmRkQlduMTJqSHlCMUxLODdEd0ltOTMKVWxhMEpCaitLd2tzMkxOeDdzR0hHSjVKUFJtMmdiamRnK3hHU0djOXlRaXM0bzFnR3ZPOUsrWTRxTmdPRFFXSApqdVZUT0FmZVRqSGdhRWlyNUZMRUdaNndZcnBDNDhJclJwVVBkbjRoTXF2M213QmphMUVpRTZlaE9zaC9yT2ZsCjhnUnBRVkJHbmZEOXhjN3o3L1F4RFg1ejloRDUzZkFxK1FWUGtMTTFScXpNZlRSZlhxNXNKUHBoenZjeHlUWTIKVHAvOWdxbkYvK21vaHR1WVVRZThwMEkrMkM3dGQyWWxPbEhuU3ZkcjJ0SklFMWVzRkE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBeVJGbC8xbEI4T1FKZ2RYNWc5WHcxWG5jZlBEazlDd1hJRXE1emhSdisrUDR5Q0w3CkQxOW9OaHJ4dmZDNFZ5L09USit4WjFzaThNTTNwbXFqQjREZVdhS1k2QUJVU3Jsc1U0MXJ1K0J3Qld3d0FhS3QKc054bS9jS3llRXRQVytGWE42RkQxZS9GdkhpVHhVeWF1WHU3aWJnNVBwaElZMCtFMUU4eHBZVlE2ZDMxd0RLZQppRDBGdzZSUEdUWHRoeWJoL3pnVUcvYVVyU25kcjdSZTBkby8yOWR4d1M4ZFhkelN6Vkp5ZGNPMUxUM1B4aStmCmZ5cU44RGxJeWdTWTB1TWdQMVVpL3p5ekNIOS9odjR2M3lzUWRQU3lMWjJObmxQaGNTbkxZY05WSzZPRzRFTHEKbkthNUFWU3VDZWkvcGFVNWVOWm1wWTI1clZiWTEyRW9XWDRKOVFJREFRQUJBb0lCQUhENUpKMG5zcVErREpWMQorNTgybXYzblhacVo0NXJLRnloZXRDMTFVRUp3b2YwWm9vVE5yNGtPOUpzcllQZ2o5VDhhVHh4cW1hUTdtODRXCmhvbVZ2OWtQUWVqQkhwdHB1ZExzV0FjVklQdjdBTEk2akdzZU80UURYREc0NXRzU0x1WWo1ekNTYTBEdFl1SHQKbTlYeURycGYxMUl6VUNEMHNnWjBoTW1pc1liazNxSzRCRWlCRm5sMXpaNnpGczB5bW5lcmIrdmNaMmlKWlZKQwovYldERHBzMThTcTFkNjVXWkRUQVBESjZuQ2ttZlkrMFRNaHJoZVVvU2RWNXl2Mjd0WjdCS2poVkp2NEZVTXRYCnM0anR6UDBMR05NMm1VSHlJMmRVSFk0WHBuRlp1OTk4dFBkT1lKdjhlT2lma3ZBYTB0Q2FpaEYraWJtZThRM3oKeVBVWktPRUNnWUVBK2dVTkxBYXRYUEdiQmwvN3RGTE9VN3J3ZnZWMVg0R0NUWkxyWGh0ZHVXSjNyejhCWXdIbwp5OElleGxqZzQ0VGN2Q0psTjJxdFdVbDhpQlRBN2xvd2VBYXF1blBIeXJERkJOT2pQT2lRVDZkNXl0OGQrdVoyClNXbFdCZGZuNDFxZkMrM3R3YW90ZUIreFFvWWczOUlDZU02ZEpiNEx6QnFqK2VhQmo3VGtaamtDZ1lFQXplQ1oKbmo3ZHljVWNPS1BCRWtOcVdVTHNzeHk2U2tjYVpKMTBNWUcxQmxuQWNyTXY0eXQ1cHNvMXQzUXR5ek5wU0VRMAp4WCtWMGpYeHBSM1Q0T3laaGNVa2tMUTZoVjBFVTV4amgyYmI1VVpkSzRUc2ljdDN6ZnV4SytUYVVsSzNVUVlBCjF6NitXOUdNNnFDeEJvdXM2UWlEZUdWNFJSbVBNYzM5RDJvdUlaMENnWUJKMklZY2M5aHV0Rm5pbHlVRUROVXUKTysvZXYrY3NEMWVLVEJQQzF2aG0waU9obWsxeGpkWFJvUU52K0Q1c1ZRalJkS05LcS9LSmhjQUhmWUhDcE5iZQpETElPb2pYbnN6QnVEWWlRYSt4NzZtNTVzazdybFhob2xIazQvcDFhTktEM2FBNHFmbll4bmNMQWNGZHpaRmNnCmtYNHU4S1FSWStqODNjMTB3YXdZV1FLQmdIVnZRMUV2M1FWRjZUWXB1bXozaFlkTEZJZUt3SHkwR0VzQ1FaQW4KdzZIT3VtYnk3VUx3dnFDMFFYWFcwSFJUYkIrcndFcFQzNXNiNkMzZVdNaTVUVTB2eWc2OHI1RDVJUW1zY0YyQQovN2ZGcmttdnRkbFg5WXZLb1NJL2xZVlY1M2xSS2xPZFY2bzZXdDVQTVF0aHl3eldMV0FVeVVqcExuUWpZdEZUCktPTkZBb0dBWUI4TDEvTUUwTHdrYkFTSVh0UmtZWHJ3VG1YK0l4WGhDNTU2ZkEybXRhRXlwSk83c0xFYk81RjYKbWZzSkUwMW1DcXhuY2RHUGcycjRQZHF1MVExRGlzOW94QUI4T215RUdTa3Q5MHFCVXp6V1VhZUhJZmFzVTIwRAoybjdOMFlDTk0vQ3c2RnU0bldKNjRDZkxFMVNBT2ZDZFYwZTlWMnVNZGVyNnpOYUN3ZDg9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
`

func NewNodeController() NodeController {
	klog.Info("NewNodeController")
	//config, err := clientcmd.RESTConfigFromKubeConfig([]byte(s))
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	return NodeController{clientSet: clientSet}
}

func (n *NodeController) Run(ctx context.Context) error {
	klog.Info("开始运行run")
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			err := n.Update(ctx)
			if err != nil {
				klog.Infoln(err)
			}
		case <-ctx.Done():
			klog.Info("程序退出")
			return nil
		}
	}
}

func (n *NodeController) Update(ctx context.Context) error {
	klog.Info("开始获取节点信息")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Error("list nodes failed")
	}
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
		details, err := eks.DescribeNodeDetails(consts.ClusterId, node.GetName())
		if err != nil {
			return err
		}
		flag, err := UpdateNode(&node, details.Data)
		if err != nil {
			return err
		}
		if flag {
			clientSet.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
		}
	}
	return nil
}

// UpdateNode 设置节点的污点
func UpdateNode(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) (bool, error) {
	klog.Info("更新节点")
	if detail == nil || node == nil {
		return false, errors.New("invalid node")
	}
	labelFlag := UpdateNodeLabels(node, detail)
	taintFlag := UpdateNodeTaints(node, detail)
	return labelFlag && taintFlag, nil
}

// UpdateNodeLabels 更新节点标签
func UpdateNodeLabels(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) bool {
	klog.Info("更新节点的标签")
	if len(detail.Labels) == 0 {
		return false
	}
	labels := make(map[string]string)
	if len(node.Labels) > 0 {
		for key, value := range node.Labels {
			labels[key] = value
		}
	}
	for i := 0; i < len(detail.Labels); i++ {
		label := detail.Labels[i]
		labels[label.Key] = label.Value
	}
	node.Labels = labels
	return true
}

// UpdateNodeTaints 修改节点的污点
func UpdateNodeTaints(node *v1.Node, detail *commoneks.DescribeEKSNodeResponseData) bool {
	klog.Info("更新节点的污点")
	taints := make([]v1.Taint, 0, 0)
	taintMap := make(map[string]v1.Taint)
	if len(detail.Taints) == 0 {
		return false
	}
	for i := 0; i < len(node.Spec.Taints); i++ {
		taint := node.Spec.Taints[i]
		taintMap[taint.Key] = v1.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: taint.Effect,
		}
	}
	for i := 0; i < len(detail.Taints); i++ {
		taint := detail.Taints[i]
		taintMap[taint.Key] = v1.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: v1.TaintEffect(taint.Effect),
		}
	}
	for _, value := range taintMap {
		taints = append(taints, value)
	}
	node.Spec.Taints = taints
	return true
}
