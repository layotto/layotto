module github.com/layotto/layotto

go 1.14

require (
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.0
	github.com/layotto/layotto/components v0.0.0-20210603045430-66065fa0b67f
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasthttp v1.26.0
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/examples v0.0.0-20210526223527-2de42fcbbce3 // indirect
	google.golang.org/protobuf v1.26.0
	mosn.io/api v0.0.0-20210414070543-8a0686b03540
	mosn.io/mosn v0.22.1-0.20210425073346-b6880db4669c
	mosn.io/pkg v0.0.0-20210401090620-f0e0d1a3efce
)

replace github.com/layotto/layotto/components v0.0.0-20210603045430-66065fa0b67f => ./components
