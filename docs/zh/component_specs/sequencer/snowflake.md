# Snowflake

## 配置项说明

示例：configs/runtime_config.json

| 字段          | 必填 | 说明                                                         |
| ------------- | ---- | ------------------------------------------------------------ |
| mysqlHost     | Y    | mysql服务器地址，比如localhost:3306                                           |
| userName      | Y    | mysql用户名                                                  |
| password      | Y    | mysql密码                                                    |
| databaseName  | Y    | mysql数据库名                                                |
| tableName     | N    | mysql表名                                                    |
| boostPower    | N    | RingBuffer扩容参数，默认3，扩容后buffersize = 8192 << 3=65536。id需求量大时可适当调大 |
| paddingFactor | N    | 少于总容量的该百分比向RingBuffer中填充UID，取值为（0，100），默认为50。id需求量大时可适当调大 |
| timeBits      | N    | 时间戳所占位数大小。默认为28                                 |
| workerBits    | N    | 机器id所占位数大小。默认为22                                 |
| seqBits       | N    | 序列号所占位数大小。默认为13                                 |
| startTime     | N    | 时间基点。默认为“2022-01-01”                                 |

## 整体设计

雪花算法生成id的整体设计如下图：

![img.jpg](../../../img/sequencer/snowflake/snowflake_id.jpg)

## 怎么启动 mysql

如果想启动snowflake的demo，需要先用Docker启动一个mysql命令：

>如果3306端口被其他服务占用，需要先退出其他服务

```shell 
docker pull mysql:latest
docker run --name snowflake -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
docker exec -it snowflake bash
mysql -uroot -p123456
```

需要在mysql中新建一个数据库：

```mysql
CREATE DATABASE layotto_sequencer;
```



## 启动 layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

>如果 build 报错，可以在项目根目录执行 `go mod vendor`

编译成功后执行:

````shell
./layotto start -c ../../configs/config_snowflake.json
````

## 运行 Demo

````shell
cd ${project_path}/demo/sequencer/common/
 go build -o client
 ./client -s "sequencer_demo"
````