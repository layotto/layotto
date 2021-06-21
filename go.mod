module mosn.io/layotto

go 1.14

require (
	github.com/dapr/components-contrib v1.2.0
	github.com/dapr/kit v0.0.1
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.26.0
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/examples v0.0.0-20210526223527-2de42fcbbce3 // indirect
	google.golang.org/protobuf v1.26.0
	mosn.io/api v0.0.0-20210414070543-8a0686b03540
	mosn.io/layotto/components v0.0.0-20210617023402-38c9d79ebb37
	mosn.io/layotto/spec v0.0.0-20210617023402-38c9d79ebb37
	mosn.io/mosn v0.23.0
	mosn.io/pkg v0.0.0-20210604065522-6e8f5a087814
	mosn.io/proxy-wasm-go-host v0.0.0-20210312032409-2334f9cf62ec
)

replace (
	mosn.io/layotto/components => ./components
	mosn.io/layotto/spec => ./spec
)
