package provider

import (
	"context"
	"fmt"
	"github.com/capitalonline/eks-cloud-controller-manager/pkg/common/consts"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"testing"
)

var s = `

`

func TestLoadBalancer_GetLoadBalancer(t *testing.T) {
	loadBalancer := LoadBalancer{}
	service := v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "dyl-slb",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec:   v1.ServiceSpec{},
		Status: v1.ServiceStatus{},
	}
	loadBalancer.GetLoadBalancer(context.Background(), "", &service)
}

func TestLoadBalancer_GetLoadBalancerName(t *testing.T) {
	loadBalancer := LoadBalancer{}
	service := v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "dyl-slb",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec:   v1.ServiceSpec{},
		Status: v1.ServiceStatus{},
	}
	loadBalancer.GetLoadBalancerName(context.Background(), "", &service)
}

func TestLoadBalancer_EnsureLoadBalancerDeleted(t *testing.T) {
	loadBalancer := LoadBalancer{}

	service := v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "dyl-slb",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec:   v1.ServiceSpec{},
		Status: v1.ServiceStatus{},
	}
	loadBalancer.EnsureLoadBalancerDeleted(context.Background(), "", &service)
}

func TestLoadBalancer_UpdateLoadBalancer(t *testing.T) {
	loadBalancer := LoadBalancer{}

	service := v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "dyl-slb",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec:   v1.ServiceSpec{},
		Status: v1.ServiceStatus{},
	}
	node := make([]*v1.Node, 0, 1)
	loadBalancer.UpdateLoadBalancer(context.Background(), "", &service, node)
}

func TestLoadBalancer_EnsureLoadBalancer(t *testing.T) {

	fmt.Println(consts.VpcID)
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(s))

	if err != nil {
		log.Fatalf("newCloud:: Failed to create kubernetes config: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	loadBalancer := LoadBalancer{
		clientSet: clientSet,
	}

	service := v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "dyl-slb",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations: map[string]string{
				AnnotationLbType:      "4",
				AnnotationLbSpec:      LBSpecStandard,
				AnnotationLbEip:       "1",
				AnnotationLbBandwidth: "50",
				AnnotationLbProtocol:  "TCP",
				AnnotationLbAlgorithm: "rr",
			},
			OwnerReferences: nil,
			Finalizers:      nil,
			ManagedFields:   nil,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:        "ingress-30443",
					Protocol:    "TCP",
					AppProtocol: nil,
					Port:        0,
					TargetPort:  intstr.IntOrString{},
					NodePort:    30443,
				},
			},
			Selector:                      nil,
			ClusterIP:                     "",
			ClusterIPs:                    nil,
			Type:                          "",
			ExternalIPs:                   nil,
			SessionAffinity:               v1.ServiceAffinityNone,
			LoadBalancerIP:                "",
			LoadBalancerSourceRanges:      nil,
			ExternalName:                  "",
			ExternalTrafficPolicy:         "",
			HealthCheckNodePort:           0,
			PublishNotReadyAddresses:      false,
			SessionAffinityConfig:         nil,
			IPFamilies:                    nil,
			IPFamilyPolicy:                nil,
			AllocateLoadBalancerNodePorts: nil,
			LoadBalancerClass:             nil,
			InternalTrafficPolicy:         nil,
		},
		Status: v1.ServiceStatus{},
	}
	nodeList, err := loadBalancer.clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	nodes := make([]*v1.Node, 0)
	for i := 0; i < len(nodeList.Items); i++ {
		node := nodeList.Items[i]
		nodes = append(nodes, &node)
	}
	loadBalancer.EnsureLoadBalancer(context.Background(), "", &service, nodes)
}
