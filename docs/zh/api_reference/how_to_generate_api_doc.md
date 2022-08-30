# 如何基于proto文件生成代码、文档

```shell
make proto
```

Then you get:
- `.pb.go` code
- API reference docs
- updated API reference list
- quickstart document (both chinese and english)
- updated sidebar (The tool will add the generated quickstart doc into the sidebar of https://mosn.io/layotto )
- updated CI (The tool will add the generated quickstart doc into the CI script `etc/script/test-quickstart.sh`)

That's all :)

## 如何把 proto 文件编译成`.pb.go`代码
<!-- tabs:start -->
### **Make 命令生成(推荐)**
本地启动 docker 后，在项目根目录下运行：

```bash
make proto.code
```

该命令会用 docker 启动 protoc，生成`.pb.go`代码。

这种方式更方便，开发者不需要修改本地 protoc 版本，省去了很多烦恼。

### **手动安装工具**
1. Install protoc version: [v3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3)

2. Install protoc-gen-go v1.28 and protoc-gen-go-grpc v1.2

3. Generate gRPC `.pb.go` code

```bash
cd spec/proto/runtime/v1
protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative *.proto
```
<!-- tabs:end -->
## 如何基于proto文件生成接口文档

我们可以用[protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) 和docker来生成接口文档，相关命令如下：

<!-- tabs:start -->
### **Make 命令生成(推荐)**
本地启动 docker 后，在项目根目录下运行：

```bash
make proto.doc
```

该命令会用 docker 启动 protoc-gen-doc，生成文档

### **用 docker 启动 protoc-gen-doc**
`make proto.doc` invokes the script `etc/script/generate-doc.sh`, which uses docker to run protoc-gen-doc.

You can check `etc/script/generate-doc.sh` for more details.

<!-- tabs:end -->