# Use Sequencer API to generate distributed unique and self-incrementing id
## What is Sequencer API
The Sequencer API is used to generate distributed unique, self-incrementing IDs.

The Sequencer API supports the declaration of demand for self-increment, including trend increase (WEAK) and strictly global increment (STRONG)

## Quick start

This example shows how to call Etcd through Layotto to generate a distributed unique, self-increasing id.

The architecture of this example is shown in the figure below, and the processes started are: Etcd, Layotto, and client programs

![img.png](../../../img/sequencer/etcd/img.png)
### step 1. Deploy etcd and Layotto
<!-- tabs:start -->
#### **With Docker Compose**
You can start etcd and Layotto with docker-compose

```bash
cd docker/layotto-etcd
# Start etcd and layotto with docker-compose
docker-compose up -d
```
#### **Compile locally (not for Windows)**
You can run etcd with Docker, then compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker-compose
#### step 1.1 Deploy the storage system (Etcd)

For the deployment of etcd, please refer to etcd's [Official Document](https://etcd.io/docs/v3.5/quickstart/)

Brief description:

Visit https://github.com/etcd-io/etcd/releases to download etcd of the corresponding operating system (docker is also available)

Once the download is finished,execute the command to start:

```shell @background
./etcd
```

The default listening address is `localhost:2379`
#### step 1.2 Run Layotto

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

Once finished, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/runtime_config.json
```
<!-- tabs:end -->

### step 2. Run the client program and call Layotto to generate a unique id

```shell
 cd ${project_path}/demo/sequencer/common/
 go build -o client
 ./client -s "sequencer_demo"
```

If the following information is printed, the demo is successful:

```bash
Try to get next id.Key:key666 
Next id:next_id:1  
Next id:next_id:2  
Next id:next_id:3  
Next id:next_id:4  
Next id:next_id:5  
Next id:next_id:6  
Next id:next_id:7  
Next id:next_id:8  
Next id:next_id:9  
Next id:next_id:10  
Demo success!
```
### step 3. Stop containers and release resources
<!-- tabs:start -->
#### **Docker Compose**
If you started etcd and Layotto with docker-compose, you can shut them down as follows:

```bash
cd ${project_path}/docker/layotto-etcd
docker-compose stop
```

#### **Destroy the etcd container**
If you started etcd with Docker, you can destroy the etcd container as follows:
```shell
docker rm -f etcd
```
<!-- tabs:end -->

### Next step
#### What does this client program do?
The demo client program uses the golang version SDK provided by Layotto, calls the Layotto Sequencer API, and generates a distributed unique, self-increasing id.

The sdk is located in the `sdk` directory, and users can call the API provided by Layotto through the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.

#### Want to learn more about Sequencer API?
What does the Sequencer API do, what problems it solves, and in what scenarios should I use it?

If you are confused and want to know more details about the use of Sequencer API, you can read [Sequencer API Usage Document](en/api_reference/sequencer/reference)

#### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
