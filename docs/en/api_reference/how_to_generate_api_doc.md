# How to generate api documents

We can use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and docker to generate api documents,the command is as follows:  
(Run in layotto directory)

```
docker run --rm \
-v  $(pwd)/docs/en/api_reference:/out \
-v  $(pwd)/spec/proto/runtime/v1:/protos \
pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,runtime_v1.md runtime.proto
```

and 

```shell
docker run --rm \
-v  $(pwd)/docs/en/api_reference:/out \
-v  $(pwd)/spec/proto/runtime/v1:/protos \
pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,appcallback_v1.md appcallback.proto
```