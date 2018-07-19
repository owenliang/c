package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
)

// 使用租约实现kv自动过期
func demo5() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		leaseGrantResp *clientv3.LeaseGrantResponse
		putResp *clientv3.PutResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse	// 只读的channel
		keepResp *clientv3.LeaseKeepAliveResponse
		lease clientv3.Lease
		leaseId clientv3.LeaseID
		ctx context.Context
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

	// 用于读写etcd键值对
	kv = clientv3.NewKV(client)

	// 用于管理lease租约
	lease = clientv3.NewLease(client)

	// 创建10秒租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 租约的ID
	leaseId = leaseGrantResp.ID
	fmt.Println("租约ID:", leaseId)

	// put一个kv
	if putResp, err = kv.Put(context.TODO(), "/cron/job5", "echo hello;", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	// put的结果
	fmt.Println("写入版本:", putResp.Header.Revision)

	// 我们持续续租5秒后停止程序
	ctx, _ = context.WithTimeout(context.TODO(), 5 * time.Second)

	// 持续给lease续约
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println("自动续租失败", err)
		return
	}

	// 消费自动续租的应答, 直到租约被取消或者出错
	for {
		select {
		case keepResp = <-keepRespChan:
			if keepResp == nil {
				fmt.Println("终止租约")
				goto END
			} else {
				fmt.Println("续租成功:", keepResp.ID)
			}
		}
	}
END:
}

func main() {
	demo5()
}