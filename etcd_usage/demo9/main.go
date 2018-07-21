package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
)

// 事务操作实现乐观锁
func demo9() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		leaseGrantResp *clientv3.LeaseGrantResponse
		lease clientv3.Lease
		leaseId clientv3.LeaseID
		txn clientv3.Txn
		ctx context.Context
		cancelFunc context.CancelFunc
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse	// 只读的channel
		keepResp *clientv3.LeaseKeepAliveResponse
		txnResp *clientv3.TxnResponse
		rangeResp *etcdserverpb.RangeResponse
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

	/* 上锁部分 */

	// 用于管理lease租约
	lease = clientv3.NewLease(client)

	// 创建5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// 准备一个用于终止续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 租约的ID
	leaseId = leaseGrantResp.ID
	fmt.Println("租约ID:", leaseId)

	// 确保函数末尾停止自动续租协程
	defer cancelFunc()
	// 确保函数末尾释放Lease
	defer lease.Revoke(context.TODO(), leaseId)

	// 自动给lease续约
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println("启动自动续租:", err)
		return
	}

	// 启动一个协程处理自动续租的应答
	go func() {
		// 消费自动续租的应答, 直到租约被取消或者出错
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("停止续租")
					goto END
				} else {
					fmt.Println("续租成功:", keepResp.ID)
				}
			}
		}
		END:
	}()

	// 用于读写etcd键值对
	kv = clientv3.NewKV(client)

	// 创建事务
	txn = kv.Txn(context.TODO())

	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/lock/job9"), "=", 0)).
		// 则抢占建立
		Then(clientv3.OpPut("/lock/job9", "I am demo9", clientv3.WithLease(leaseId))).
		// 否则乐观锁抢占失败, 获取其值
		Else(clientv3.OpGet("/lock/job9"))

	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println("事务失败:", err)
		return
	}

	// If条件未成立, 抢锁失败
	if !txnResp.Succeeded {
		rangeResp = txnResp.Responses[0].GetResponseRange()
		fmt.Println("锁被占用:", string(rangeResp.Kvs[0].Value))
		return
	}

	/* 锁内时间, 执行业务处理 */

	fmt.Println("成功获锁")
	time.Sleep(5 * time.Second)

	/* defer会自动释放lease, 同时停止续租 */
}

func main() {
	demo9()

	// 等1秒, 为了演示自动续租终止
	time.Sleep(1 * time.Second)
}