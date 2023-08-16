package provider

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

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
