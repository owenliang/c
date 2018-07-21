package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
)

// get读取一条记录
func demo3() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		getResp *clientv3.GetResponse
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

	// 读取一条记录
	if getResp, err = kv.Get(context.TODO(), "/cron/job1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs)
	}
}

func main() {
	demo3()
}