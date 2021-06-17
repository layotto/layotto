package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"strconv"
	"sync"
	"time"
)

const (
	storeName = "apollo"
	appid     = "testApplication_yang"
	group     = "application"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()

	// Belows are CRUD examples
	// 1. set
	testSet(ctx, cli)
	// 2. get after set
	// Since configuration data in apollo cache is eventual-consistent,we need to sleep a while before querying new data
	time.Sleep(time.Second * 2)
	testGet(ctx, cli)

	// 3. delete
	testDelete(ctx, cli)

	// 4. get after delete
	//sleep a while before querying deleted data
	time.Sleep(time.Second * 2)
	testGet(ctx, cli)

	// 5. show how to use subscribe API
	// with sdk
	//testSubscribeWithSDK(ctx, cli)

	// besides sdk,u can also call our server with grpc
	testSubscribeWithGrpc(ctx)
}

func testSubscribeWithSDK(ctx context.Context, cli client.Client) {
	var wg sync.WaitGroup
	wg.Add(1)
	// 1. subscribe
	ch := cli.SubscribeConfiguration(ctx, &client.ConfigurationRequestItem{
		StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"haha", "heihei"},
	})
	// 2. client loop receiving changes in another gorountine
	go func() {
		for resp := range ch {
			if resp.Err != nil {
				fmt.Println("subscribe failed", resp.Err)
				continue
			}
			marshal, err := json.Marshal(resp.Item)
			if err != nil {
				fmt.Println("marshal resp.Item failed.")
				return
			}
			fmt.Printf("receive subscribe resp: %+v\n", string(marshal))
		}
		fmt.Println("subscribe channel has been closed.")
	}()

	// 3. loop setting configuration in another gorountine
	go func() {
		idx := 0
		for {
			time.Sleep(1 * time.Second)
			idx++
			newItmes := make([]*client.ConfigurationItem, 0, 10)
			newItmes = append(newItmes, &client.ConfigurationItem{Group: "application", Label: "prod", Key: "heihei", Content: "heihei" + strconv.Itoa(idx)})
			fmt.Println("write start")
			cli.SaveConfiguration(ctx, &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid, Items: newItmes})
		}
	}()
	wg.Wait()
}

func testSubscribeWithGrpc(ctx context.Context) {
	// 1. connect with grpc
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := runtimev1pb.NewRuntimeClient(conn)

	// get client for subscribe
	var wg sync.WaitGroup
	wg.Add(1)
	cli, err := c.SubscribeConfiguration(ctx)
	// client receive changes
	go func() {
		for {
			data, err := cli.Recv()
			if err != nil {
				fmt.Println("subscribe failed", err)
				continue
			}
			fmt.Println("receive subscribe resp", data)

		}
	}()
	// client send subscribe request
	go func() {
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"haha", "key1"}})
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"heihei"}})
	}()

	// loop write in another gorountine
	go func() {
		idx := 0
		for {
			time.Sleep(1 * time.Second)
			idx++
			newItmes := make([]*runtimev1pb.ConfigurationItem, 0, 10)
			newItmes = append(newItmes, &runtimev1pb.ConfigurationItem{Group: "application", Label: "prod", Key: "heihei", Content: "heihei" + strconv.Itoa(idx)})
			fmt.Println("write start")
			c.SaveConfiguration(ctx, &runtimev1pb.SaveConfigurationRequest{StoreName: storeName, AppId: appid, Items: newItmes})
		}
	}()
	wg.Wait()
}

func testSet(ctx context.Context, cli client.Client) {
	item1 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "key1", Content: "value1", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "print",
	}}
	item2 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "haha", Content: "heihei", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "haha",
	}}
	saveRequest := &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)
	if cli.SaveConfiguration(ctx, saveRequest) != nil {
		fmt.Println("save key failed")
		return
	}
	fmt.Println("save key success")
}

func testGet(ctx context.Context, cli client.Client) {
	getRequest := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	items, err := cli.GetConfiguration(ctx, getRequest)
	//validate
	if err != nil {
		fmt.Printf("get configuration failed %+v \n", err)
	}
	for _, item := range items {
		fmt.Printf("get configuration after save, %+v \n", item)
	}
}

func testDelete(ctx context.Context, cli client.Client) {
	request := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	if cli.DeleteConfiguration(ctx, request) != nil {
		fmt.Println("delete key failed")
		return
	}
	fmt.Printf("delete keys success\n")
}
