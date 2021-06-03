# ʹ��Layotto Actuator���н�������Ԫ���ݲ�ѯ

��ʾ��չʾ�����ͨ��Layotto Actuator��Http API���н�������Ԫ���ݲ�ѯ

## ʲô��Layotto Actuator

�����������У���Ҫ��Ӧ�ó����״̬���м�أ���Layotto�Ѿ�������һ����ع��ܣ�����Actuator�� ʹ��Layotto Actuator���԰������غ͹���Layotto��Layotto�����Ӧ�ã����罡����顢��ѯ����ʱԪ���ݵȡ�
���е���Щ���Կ���ͨ��HTTP�ӿ������ʡ�

## ���ٿ�ʼ

### ����Layotto server ��

����Ŀ�������ص����غ��л�����Ŀ¼�����룺

```bash
cd ${projectpath}/cmd/layotto
go build
```

��ɺ�Ŀ¼�»�����layotto�ļ�����������

```bash
./layotto start -c ../../configs/config_apollo_health.json
```

### ���ʽ������ӿ�

���� /actuator/health/liveness

```bash
curl http://127.0.0.1:34999/actuator/health/liveness
```

���أ�

```json
{
  "components": {
    "apollo": {
      "status": "UP"
    },
    "runtime_startup": {
      "status": "UP",
      "details": {
        "reason": ""
      }
    }
  },
  "status": "UP"
}
```

����"status": "UP"����״̬��������ʱ���ص�Http״̬����200��

### ��ѯԪ����

���� /actuator/info

```shell
curl http://127.0.0.1:34999/actuator/info
```

���أ�

```json
{
  "app": {
    "name": "Layotto",
    "version": "0.1.0",
    "compiled": "2021-05-20T14:32:40.522057+08:00"
  }
}
```

### ģ�����ô���ĳ���

���Layotto���ô������������������ṩ����ͨ��������鹦�ܿ��Լ�ʱ���֡�

���ǿ���ģ��һ�����ô���ĳ�����ʹ��һ������������ļ�����Layotto:

```shell
./layotto start -c ../../configs/wrong/config_apollo_health.json
```

�������ļ������������˷���apollo��Ҫ��open_api_token��

���ʽ������ӿڣ�ע���������õĶ˿���34888������һ�������в�һ������

```shell
curl http://127.0.0.1:34888/actuator/health/liveness
```

���أ�

```json
{
  "components": {
    "apollo": {
      "status": "DOWN",
      "details": {
        "reason": "configuration illegal:no open_api_token"
      }
    },
    "runtime_startup": {
      "status": "DOWN",
      "details": {
        "reason": "configuration illegal:no open_api_token"
      }
    }
  },
  "status": "DOWN"
}
```

json��"status": "DOWN"����ǰ״̬����������ʱ���ص�Http״̬����503��


## ��һ��

### ���ɽ�Kubernetes�������

Layotto�����ṩ��/actuator/health/readiness��/actuator/health/liveness �����������ӿڣ���ӦKubernetes������鹦����Readiness��Liveness�������塣

��ˣ������Բ���[Kubernetes���ĵ�](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) �����������ӿڼ��ɽ�Kubernetes������顣

### Ϊ���������ӽ�������Ԫ���ݲ�ѯ����

�����ʵ�����Լ���Layotto���������Ϊ����ӽ���������������Բο�apollo�����ʵ�֣��ļ���components/configstores/apollo/indicator.go����ʵ��info.Indicator�ӿڣ�������ע���Actuator���ɡ�

### �˽�Actuatorʵ��ԭ��

�������ʵ��ԭ�����Ȥ����������Actuator��չһЩ���ܣ������Ķ�[Actuator������ĵ�](../../design/actuator/actuator-design-doc.md)