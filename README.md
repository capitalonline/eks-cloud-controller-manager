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
    service.beta.kubernetes.io/cds-load-balancer-protocol: TCP
    service.beta.kubernetes.io/cds-load-balancer-types: "4"
    service.beta.kubernetes.io/cds-load-balancer-specification: standard
    service.beta.kubernetes.io/cds-load-balancer-bandwidth: "50"
    service.beta.kubernetes.io/cds-load-balancer-eip: "1"
    service.beta.kubernetes.io/cds-load-balancer-algorithm: rr
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

| key                                                        | val            | 说明                                                         |
| ---------------------------------------------------------- | -------------- | ------------------------------------------------------------ |
| service.beta.kubernetes.io/cds-load-balancer-protocol      | TCP/UDP        | 协议类型，只有TCP、UDP                                       |
| service.beta.kubernetes.io/cds-load-balancer-types         | 4              | 网络模型四层网络                                             |
| service.beta.kubernetes.io/cds-load-balancer-specification | standard       | SLB规格：standard-标准型、high-高阶型、super-超强型、extreme-至强型。客户根据实际产品分配填写 |
| service.beta.kubernetes.io/cds-load-balancer-bandwidth     | 50             | 共享带宽大小，单位M，10的倍数，范围10~1000                   |
| service.beta.kubernetes.io/cds-load-balancer-eip           | 1              | 创建的eip数量，目前只传1                                     |
| service.beta.kubernetes.io/cds-load-balancer-algorithm     | rr/wrr/conhash | 监听轮询协议：rr-轮询、wrr-加权轮询、conhash-一致性哈希      |
| service.beta.kubernetes.io/cds-load-balancer-subject-id    |                | 测试金id，客户试用时填写，具体值联系管理员                   |

**注意：**

1. 在创建LoadBalancer的Service前，确保用户具有SLB产品的使用权限和配额
2. 创建好LoadBalancer类型后，用户可在首云gic页面`私有网络`->`高性能负载均衡`->`实例管理`下查询到由CCM创建的SLB实例（按需计费）以及相应的监听规则，同时包含相关共享带宽和EIP
3. 当删除了Service-LoadBalancer后，为防止用户SLB下其他的非CCM自动创建监听规则被勿删，CCM不会删除SLB实例，而是清除其相关监听规格。用户若不需要此SLB，请自行释放资源以免不必要的计费
