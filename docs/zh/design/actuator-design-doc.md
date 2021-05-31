# Actuator����ĵ�
# һ����Ʒ���
## 1.1. ����

- �������

ͨ��Actuator�ӿڿ���ͳһ��ȡ��Layotto�ڲ���������Լ�ҵ��Ӧ�õĽ���״̬

- �鿴����ʱԪ����

ͨ��Actuator�ӿڿ���ͳһ��ȡ��Layotto�Լ���Ԫ������Ϣ������汾��git��Ϣ�����Լ�ҵ��Ӧ�õ�Ԫ������Ϣ������ͨ���������Ķ��ĵ��������б�����Ӧ�ð汾��Ϣ��

- ֧�ּ��ɽ���Դ������ʩ��������
    - ���Լ��ɽ�k8s�������
    - ���Լ��ɽ����ϵͳ������Prometheus+Grafana
    - ������Ҫ��ע�����Ŀ��Ի��ڽ���������޳��ڵ�
    - �������Ի��ڴ˽ӿ���dashboard��Ŀ����GUI����,�Ա��Ų����⡣
    
- ������Spring Boot Actuator�Ĺ��ܣ�δ���и��������ռ䣺Monitoring, Metrics, Auditing, and more.

## 1.2. ����

**Q: ��ֵ��ɶ���������ӿڿ�������˭�ã�**

1. �������Ų����⣬ֱ�ӵ��ӿڲ�ѯ����ʱ��Ϣ����������dashboardҳ��/GUI����

2. �����ϵͳ����أ�

3. ��������ʩ���Զ�����ά�����粿��ϵͳ���ڽ���������жϲ�����ȣ�ֹͣ������������𣻱���ע�����Ļ��ڽ�������޳��쳣�ڵ㣻����k8s���ڽ������kill���������´�������


**Q: ���񷵻ظ�״̬����У�û��Ҫ��������ʱ��Ϣ�������������ʱ��ϸ��Ϣ��˭�ã�**

1. �������Ի��ڴ˽ӿ���dashboardҳ�����GUI����,�Ա��Ų����⣻

������spring boot��������spring boot actuatorд�˸�spring boot admin��ҳ
�ο�[https://segmentfault.com/a/1190000017816452](https://segmentfault.com/a/1190000017816452)

2. ���ɼ��ϵͳ:���Խ���Prometheus+Grafana

������Spring Boot Actuator����Prometheus+Grafana
�ο�[Spring-Boot-Metrics���֮Prometheus-Grafana](https://bigjar.github.io/2018/08/19/Spring-Boot-Metrics���֮Prometheus-Grafana/)


**Q: �������ܿ����������硰���� Layotto �ڲ��ض������������**

A: ���������ز����������app����partial failure״̬���в�ȷ���ԡ�
���Ǻ������Կ������debug����������mock��ץ���İ���


**Q: �������Ľӿ�������Ȩ�޹ܿ�**

A: �Ȳ��㣬�з��������ټӸ�����


# ������Ҫ���

## 2.1. ���巽��

�ȿ���http�ӿڣ���Ϊ��Դ������ʩ�Ľ�����鹦�ܻ����϶�֧��http������k8s,prometheus)��û��֧��grpc�ġ�

Ϊ���ܹ�����MOSN�ļ�Ȩfilter��filter������Actuator����Ϊ7���filter����MOSN�ϡ�

������˵��MOSN����listener,��д��stream_filter,���filter����http����������Actuator.

Actuator�ڲ������Endpoint��������󵽴��������Actuator��ί�ж�Ӧ��Endpoint���д���Endpoint֧�ְ�����չ��ע���Actuator��

![img.png](../../../img/actuator/abstract.png)

## 2.2. Http API���

### 2.2.1. ·������

·������restful��񣬲�ͬ��Endpointע���Actuator��·����

```
/actuator/{endpoint_name}/{params}  
```

����

```
/actuator/health/liveness
```

����health��ʶEndpoint��������health��liveness�Ǵ�����Endpoint�Ĳ�����

����֧�ִ���������� /a/b/c/d�����崫������������������ÿ��Endpoint�Լ���


Ĭ��ע���·���У�

```
/actuator/health/liveness
/actuator/health/readiness
/actuator/info
```

### 2.2.2. Health Endpoint
#### /actuator/health/liveness

GET
```json
// http://localhost:8080/actuator/health/liveness
// HTTP/1.1 200 OK

{
  "status": "UP",
  "components": {
    "livenessProbe": {
      "status": "UP",
      "details":{
				 
      }
    }
  }
}
```
�����ֶ�˵����
HTTP״̬��200����ɹ�������(400���ϵ�״̬��)����ʧ��
status�ֶ������֣�
```go
var (
	// INIT means it is starting
	INIT = Status("INIT")
	// UP means it is healthy
	UP   = Status("UP")
	// DOWN means it is unhealthy
	DOWN = Status("DOWN")
)
```

#### /actuator/health/readiness

GET
```json
// http://localhost:8080/actuator/health/readiness
// HTTP/1.1 503 SERVICE UNAVAILABLE

{
  "status": "DOWN",
  "components": {
    "readinessProbe": {
      "status": "DOWN"
    }
  }
}
```
### 2.2.3. Info Endpoint

#### /actuator/info

GET
```json
// http://localhost:8080/actuator/health/liveness
// HTTP/1.1 200 OK

{
    "app" : {
        "version" : "1.0.0",
        "name" : "Layotto"
    }
}
```




**Q: ����ʱԪ����Ҫ��Щ��**

һ�ڣ�

- �汾��

�������Լ��ϣ�

- �ص�app
- ����ʱ���ò���


**Q: �Ƿ�ǿ��Ҫ�����ʵ�ֽ����ȼ��ӿڣ�**

��ʱ��ǿ��

## 2.3. �������ݵ�����ģ��

![img.png](../../../img/actuator/actuator_config.png)

����listener���ڴ���actuator��stream_filters����actuator_filter�����ڴ���actuator�����󣨼��£�

## 2.4. �ڲ��ṹ������������

![img.png](../../../img/actuator/actuator_process.png)

���ͣ�

### 2.4.1. ���󵽴�mosn��ͨ��stream filter����Layotto������Actuator

stream filter���httpЭ��ʵ����(struct)ΪDispatchFilter������http·���ַ����󡢵���Actuator:
```go

type DispatchFilter struct {
	handler api.StreamReceiverFilterHandler
}

func (dis *DispatchFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	dis.handler = handler
}

func (dis *DispatchFilter) OnDestroy() {}

func (dis *DispatchFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
}
```
Э����Actuator������δ����Ҫ����Э��Ľӿڣ�����ʵ�ָ�Э���stream filter

### 2.4.2. ����ַ���Actuator�ڲ���Endpoint

�ο�spring boot actuator����ƣ�
Actuator�����Endpoint���֧�ְ�����չ��ע��Endpoint��������ʵ��health��info Endpoint��
```go
type Actuator struct {
	endpointRegistry map[string]Endpoint
}

func (act *Actuator) GetEndpoint(name string) (endpoint Endpoint, ok bool) {
	e, ok := act.endpointRegistry[name]
	return e, ok
}

func (act *Actuator) AddEndpoint(name string, ep Endpoint) {
	act.endpointRegistry[name] = ep
}

```
������󣬸���·��������ַ�����Ӧ��Endpoint������/actuator/health/readiness��ַ���health.Endpoint

### 2.4.3. health.Endpoint������ַ���health.Indicator��ʵ��

��Ҫ�ϱ����������Ϣ�����ʵ��Indicator�ӿڡ�ע���health.Endpoint��
```go
type Indicator interface {
	Report() Health
}
```
health.Endpoint������ַ���health.Indicator��ʵ��

### 2.4.4. info.Endpoint������ַ���info.Contributor��ʵ��

��Ҫ�ϱ�����ʱ��Ϣ�����ʵ��Contributor�ӿڡ�ע���info.Endpoint��
```go
type Contributor interface {
	GetInfo() (info interface{}, err error)
}
```
info.Endpoint������ַ���info.Contributor��ʵ��

# ������ϸ���
## 3.1. ������
### 3.1.1. runtime_startup

- SetStarted���

![img.png](../../../img/actuator/set_started.png)

- SetUnhealthy���

����ʧ��:

![img.png](../../../img/actuator/img.png)

Stop��ʱ��

![img.png](../../../img/actuator/img_1.png)

### 3.1.2. apollo���

init:

![img_2.png](../../../img/actuator/img_2.png)

��ʵĿǰû����Ҫ���ĵط�����Ϊ����init��ʼ������ʧ�ܵĻ���runtime_startup��indicatorҲ�ܱ�unhealthy


