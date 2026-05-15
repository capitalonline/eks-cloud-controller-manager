# cloud-controller-manager

eks cloud-controller-manager

## Service-LoadBalancer

**使用示例**

service.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/cds-load-balancer-protocol: "TCP"
    service.beta.kubernetes.io/cds-load-balancer-types: "4"
    service.beta.kubernetes.io/cds-load-balancer-specification: "high"
    service.beta.kubernetes.io/cds-load-balancer-bandwidth: "50"
    service.beta.kubernetes.io/cds-load-balancer-network: "public"
    service.beta.kubernetes.io/cds-load-balancer-algorithm: "rr"
    service.beta.kubernetes.io/cds-load-balancer-subject-id: "00000"
  name: lb-cluster-test
  namespace: default
spec:
  externalTrafficPolicy: Cluster
  selector:
    app: nginx
  ports:
    - name: test
      protocol: TCP
      port: 10080
      targetPort: 80
  type: LoadBalancer
```

**annotations说明**

| key                                                        | val示例    | 说明                                                             |
| ---------------------------------------------------------- | -------- | -------------------------------------------------------------- |
| service.beta.kubernetes.io/cds-load-balancer-protocol      | TCP/UDP  | 协议类型，暂只支持TCP、UDP                                               |
| service.beta.kubernetes.io/cds-load-balancer-types         | 4        | 网络模型四层网络，默认4，暂只支持4层网络              |
| service.beta.kubernetes.io/cds-load-balancer-network       | public   | slb网络类型：public-公网、all-公网+私网                                    |
| service.beta.kubernetes.io/cds-load-balancer-specification | standard | SLB规格：standard-标准型、high-高阶型、super-超强型、extreme-至强型。费用详情请联系产品经理。 |
| service.beta.kubernetes.io/cds-load-balancer-billingmethod | number   | 带宽类型：number-固定带宽、flow_demand-流量按需。费用详情请联系产品经理。                 |
| service.beta.kubernetes.io/cds-load-balancer-bandwidth     | 50       | 固定带宽大小，单位M，10的倍数，范围10~1000。默认50。流量按需无视。                        |
| service.beta.kubernetes.io/cds-load-balancer-algorithm     | rr       | 监听轮询协议：rr-轮询、wrr-加权轮询、conhash-一致性哈希                            |
| service.beta.kubernetes.io/cds-load-balancer-subject-id    | 0000     | 测试金id，客户试用时填写。费用详情请联系产品经理。                                     |
| service.beta.kubernetes.io/cds-select-slb-id               | slb-uuid | 选择已有的SLB                                                       |
| service.beta.kubernetes.io/cds-select-slb-eip-addr         | 60.0.0.0 | 选择SLB公网监听地址                                                    |
| service.beta.kubernetes.io/cds-select-slb-vip-addr         | 10.0.0.0 | 选择SLB私网监听地址                                                    |

**注意：**

1. 在创建LoadBalancer的Service前，确保用户具有SLB产品的使用权限和配额
2. 创建好LoadBalancer类型后，用户可在首云gic页面`私有网络`->`高性能负载均衡`->`实例管理`下查询到由CCM创建的SLB实例（按需计费）以及相应的监听规则，同时包含相关共享带宽和EIP
3. 当删除了Service-LoadBalancer后，为防止用户SLB下其他的非CCM自动创建监听规则被勿删，CCM不会删除SLB实例，而是清除其相关监听规格。用户若不需要此SLB，请自行释放资源以免不必要的计费
