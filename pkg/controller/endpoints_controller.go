package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	AnnotationEPUpdated = "service.beta.kubernetes.io/eks-endpoints-updated"
)

// EndpointsController 监听 Endpoints 变化，触发 LoadBalancer 更新
type EndpointsController struct {
	clientSet *kubernetes.Clientset
}

func NewServiceController() (*EndpointsController, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("创建 kubernetes config 失败: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 kubernetes clientSet 失败: %v", err)
	}
	return &EndpointsController{
		clientSet: clientSet,
	}, nil
}

// Run 启动 Endpoints 监听器
func (sc *EndpointsController) Run(ctx context.Context) {
	klog.Info("start EndpointsController")

	factory := informers.NewSharedInformerFactory(sc.clientSet, time.Minute)
	endpointsInformer := factory.Core().V1().Endpoints().Informer()

	_, _ = endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldEndpoints, ok1 := oldObj.(*v1.Endpoints)
			newEndpoints, ok2 := newObj.(*v1.Endpoints)
			if !ok1 || !ok2 {
				return
			}

			// 检查 Endpoints 是否真的发生了变化
			if reflect.DeepEqual(oldEndpoints.Subsets, newEndpoints.Subsets) {
				return
			}

			// 获取对应的 Service
			service, err := sc.clientSet.CoreV1().Services(newEndpoints.Namespace).Get(ctx, newEndpoints.Name, metav1.GetOptions{})
			if err != nil {
				if !errors.IsNotFound(err) {
					klog.Errorf("get service %s/%s failed: %v", newEndpoints.Namespace, newEndpoints.Name, err)
				}
				return
			}

			// 只处理 LoadBalancer 类型的 Service
			if service.Spec.Type != v1.ServiceTypeLoadBalancer {
				return
			}

			// 只处理 externalTrafficPolicy=Local 的 Service
			if service.Spec.ExternalTrafficPolicy != v1.ServiceExternalTrafficPolicyTypeLocal {
				return
			}

			klog.Infof("detected changes in Endpoints for Service %s/%s, triggering load balancer update", service.Namespace, service.Name)

			// 触发负载均衡器更新
			if err = sc.updateLoadBalancerForService(ctx, service); err != nil {
				klog.Errorf("update Service %s/%s failed: %v", service.Namespace, service.Name, err)
			}
		},
	})

	factory.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), endpointsInformer.HasSynced) {
		klog.Fatalf("wait endpoints informer failed")
		return
	}

	klog.Info("EndpointsController Informer Running")
	<-ctx.Done()
	return
}

// updateLoadBalancerForService 触发负载均衡器更新
func (sc *EndpointsController) updateLoadBalancerForService(ctx context.Context, service *v1.Service) error {
	// 添加 annotation 触发 CCM 更新
	// CCM 会监听 Service 的变化，通过修改 annotation 可以触发更新
	if service.Annotations == nil {
		service.Annotations = make(map[string]string)
	}

	// 使用时间戳作为触发标记
	service.Annotations[AnnotationEPUpdated] = time.Now().Format(time.RFC3339)

	_, err := sc.clientSet.CoreV1().Services(service.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("update Service annotation failed: %v", err)
	}

	klog.Infof("load balancer update triggered for Service %s/%s", service.Namespace, service.Name)

	return nil
}
