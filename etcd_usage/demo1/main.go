package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
)

// 创建etcd客户端
func demo1() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},	// 集群列表
		DialTimeout: 5 * time.Second,	// 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	client = client
}

func main() {
	demo1()
}