package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// 删除kv
func demo5() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
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

	// 删除kv
	if delResp, err = kv.Delete(context.TODO(), "/cron/job1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}

	// 判断删除结果
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair = range delResp.PrevKvs {
			fmt.Println("删除了kv:", string(kvpair.Key))
		}
	}
}

func main() {
	demo5()
}