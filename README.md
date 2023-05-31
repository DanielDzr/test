## Task
        Assignment：基于nodejs的微服务，在非入侵式修改服务本身的条件下，如何全链路（包括异步服务）实现灰度发布。
	Note：请注意题设是基于微服务调用链的，例如有四个微服务A -> B ->C->D
		如果在同一个发布中，B有新版本想做灰度发布（B1）C有新版本想做灰度发布（C1）
		那整体调用应该实现为：A -> B ->C->D 
					    A ->B1->C1->D
		B，C应该有灰度的deployment并且B调用到C，而B1调用到C1
		A，D维持原先的deployment和配置

		另外请考虑有消息队列的情况下，异步过程如何实现灰度发布。


## 最后我选择的方案
我理解的链路方向：消息队列--> A --> B --> C --> D 或者 消息队列--> A --> B1 --> C1 --> D
> 微服务,无入侵，灰度发布
>> 理解中的满足task灰度发布类型有：1.AB测试 2.金丝雀发布 3.基于负载均衡器做的 4.istio 

>> 这里选择的istio

- 通过kustomize控制灰度发布的应用的deployment编写
- 通过istio控制A流量的分发（控制A把消息分发到B，B1）
- B1后端镜像构建时或者deployment填写参数时，制定post的地址为C1的svc地址

消息队列异步：
> 消息队列选用的 AWS SQS(Amazon message queuing service),个人倾向于这个选择，因为方便，快，稳定，价格便宜

1.如果是我理解中的task，消息队列是发送到A，再通过A的处理扔给B。如果，是这样的情况的的话，能实现消息队列也有10%分发到了B1。因为我server里面的post的地址是kubernetes service的地址，Istio 在这里的作用是对 Kubernetes 服务层的网络流量进行管理和控制。

2.如果消息队列是直接发送到B的话，一个可能的方法是使用一个"代理"或"网关"服务来接收所有来自消息队列的消息，然后根据一些规则（可能是基于Istio的路由规则）来决定将消息发送到哪个服务。

> 如果消息队列选择是kafka之类的话

1.如果是我理解中的task，消息通过A后走istio流量分发

2.如果消息队列是直接发送到B的话，可以创建两个KafkaConsumers，分别对应服务B和服务B1，然后根据你的灰度发布策略，调整这两个 Consumers的消费速率。


## kustomize目录
- `b`,`b1`,deployment的labels里面的version不同，用于区分灰度发布的版本
- app_a,app_b里面存放了istio virtualservice和destinationrule的yaml文件，用于控制A的流量分发
```yaml
## 部署方式
cd kustomize
kustomize build app_a | kubectl apply -f -
...
```


## server目录
- 基于Task给出的逻辑用golang编写的4个后端服务，用于测试
- - app-b,app-c 构建b1,c1镜像时，要把后端b1 post的svc地址改成c1的svc地址
- 构建镜像并推送到dockerhub
 
  **但是Dockerfile已经Image里面有相关key信息以及SQS URL就没有放到repo里展示**


