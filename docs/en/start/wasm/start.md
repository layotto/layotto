## Run business logic in Layotto using WASM

### What is WASM on Layotto?
The sidecar of service mesh and multi-runtime is a common infrastructure for the whole company, but in practice, the business system will also have its own SDK, and it will also have the difficulty of pushing users to upgrade the SDK and the problem of version fragmentation.

For example, a business system has developed an SDK in the form of a jar package for use by other business systems. Their features are not universal to the entire company, so they cannot persuade the middleware team to develop them into the company's unified sidecar.

![img_1.png](../../../img/wasm/img_1.png)

And if it becomes like this:

![img.png](../../../img/wasm/img.png)

If developers no longer develop sdk (jar package), change to develop .wasm files and support independent upgrade and deployment, there will be no pain to push the users to upgrade. 

When you want to upgrade, you can release it on the operation platform. There is no need to restart the app and sidecar.

Layotto can load the compiled WASM files automatically, and interacts with them through the API of the `proxy_abi_version_0_2_0` version.

### Quick start

#### step 1. start redis server and write test data

The example only needs a Redis server that can be used normally. As for where it is installed, there is no special restriction. It can be a virtual machine, a local machine or a server. 

Here, we run redis with docker:

Run redis container:
```shell
docker run -d --name redis-test -p 6379:6379 redis
```

execute the command`set book1 100`

```shell
docker exec -i redis-test redis-cli set book1 100
```

This command will set `book1` with 100.

The result will be:

```bash
OK
```

We can execute `get book1` to check the value of `book1`:

```shell
docker exec -i redis-test redis-cli get book1
```

The result is:

```bash
"100"
```

#### step 2. start Layotto server
Build:

```shell @if.not.exist layotto_wasmer
go build -tags wasmer -o ./layotto_wasmer ./cmd/layotto/main.go
```

Run it:
```shell @background
./layotto_wasmer start -c ./demo/faas/config.json
```

**Note: You need to modify the redis address as needed, the default address is: localhost:6379**

#### step 3. send request
```shell
curl -H 'id:id_1' 'localhost:2045?name=book1'
```

It will returns:

```bash
There are 100 inventories for book1.
```


This http request will access the wasm module in Layotto. The wasm module will call redis for logical processing.

#### step 4. Release resources

```shell
docker rm -f redis-test
```

### Note

This feature is still in the experimental stage, and the implementation of the WASM interactive API in the community is not uniform enough, so if you have any needs for this module, please post it in the issue area, we will build WASM together!