package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
)

// put写入一条记录
func demo2() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		putResp *clientv3.PutResponse
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

	// 写入一条记录
	if putResp, err = kv.Put(context.TODO(),"/cron/job1", "echo hello;", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}

func main() {
	demo2()
}