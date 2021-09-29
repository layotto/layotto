## trace跟踪

### 功能介绍

在[runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) 中，有一段关于trace的配置如下：

```json
[
  "tracing": {
    "enable": true,
    "driver": "SOFATracer",
    "config": {
      "generator": "mosntracing",
      "exporter": ["stdout"]
    }
  }
]
```
这段配置可以开启layotto的trace能力。用户可以通过配置来指定trace上报的方式，以及spanId,traceId的生成方式。

对应的调用端代码在[client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go) 中，layotto的trace打印如下：

![img.png](../../../img/trace/trace.png)


### 配置参数说明

trace配置：

| 字段名 | 字段类型 | 说明 |
|  ----  | ----  | ---- |
| enable  | boolean | 全局开关，是否开启trace|
| driver  | String | driver是代表trace的类型，mosn现有SOFATracer和skywalking，用户可以拓展|
| config  | Object | trace的拓展配置 |

trace拓展配置：

| 字段名 | 字段类型 | 说明 |
|  ----  | ----  | ---- |
| generator  | String | spanId,traceId等资源的生成方式，用户可自行拓展|
| exporter  | Array | 用户需要trace上报的方式，可自行实现和拓展|




### Trace设计和拓展

#### Span结构：

```go
type Span struct {
    StartTime     time.Time //收到请求的时间
    EndTime       time.Time //返回的时间
    traceId       string   //traceId
    spanId        string  //spanId
    parentSpanId  string  //父spanId
    tags          [xprotocol.TRACE_END]string //拓展字段，component可以将自己的信息存放到该字段
    operationName string
}
```
Span结构定义了layotto和其component之间传递的数据结构，如下图所示，component可以通过tags将自己的信息传递到layotto，layotto做
统一的trace上报：

#### generator接口：

```go
type Generator interface {
    GetTraceId(ctx context.Context) string
    GetSpanId(ctx context.Context) string
    GenerateNewContext(ctx context.Context, span api.Span) context.Context
    GetParentSpanId(ctx context.Context) string
}
```

该接口对应上面的generator配置，该接口主要用来根据收到的context生成traceId,spanId,获得父spanId以及传递给组件的context的功能，用户
可以实现自己的Generator，可以参考代码中的[OpenGenerator](../../../../diagnostics/genetator.go)的实现。

#### Exporter接口：

```go
type Exporter interface {
ExportSpan(s *Span)
}
```

exporter接口定了如何将Span的信息上报给远端，对应配置中的exporter字段，该字段是个数组，可以上报给多个服务端。可以
参照[StdoutExporter](../../../../diagnostics/exporter_iml/stdout.go)的实现,该实现将trace的信息打印到标准输出。


#### Span的上下文传递：

##### Layotto测
```go
GenerateNewContext(ctx context.Context, span api.Span) context.Context
```

GenerateNewContext用于生成新的context，我们通过mosnctx可以将该context实现上下文之间的互相传递：

```go
ctx = mosnctx.WithValue(ctx, types.ContextKeyActiveSpan, span)
```
可以参考代码中的[OpenGenerator](../../../../diagnostics/genetator.go)的实现

##### Component测

在Component测可以通过[SetExtraComponentInfo](../../../../components/trace/utils.go)塞入component的信息，
比如在接口[Hello](../../../../components/hello/helloworld/helloworld.go)执行了以下操作：

```go
	trace.SetExtraComponentInfo(ctx, fmt.Sprintf("method: %+v", "hello"))
```

trace打印的结果如下：

![img.png](../../../img/trace/trace.png)