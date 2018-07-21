package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
)

// 利用Op取代Get/Put
func demo8() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		putOp clientv3.Op
		opResp clientv3.OpResponse
		getOp clientv3.Op
		err error
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"36.111.184.221:2379"},	// 集群列表
		DialTimeout: 5 * time.Second,	// 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd键值对
	kv = clientv3.NewKV(client)

	// 生成一个put操作
	putOp = clientv3.OpPut("/cron/job8", "echo hello;")

	// 执行put操作
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入版本:", opResp.Put().Header.Revision)
	}

	// 生成一个get操作
	getOp = clientv3.OpGet("/cron/job8")

	// 执行get操作
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("数据版本:", opResp.Get().Kvs[0].ModRevision)
		fmt.Println("读到数据:", string(opResp.Get().Kvs[0].Value))
	}
}

func main() {
	demo8()
}