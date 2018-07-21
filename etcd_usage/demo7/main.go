package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// 使用watch监听目录变化
func demo7() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watcher clientv3.Watcher
		curVal string
		watchStartRev int64
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		event *clientv3.Event
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

	// 启动一个协程, 定时的更新与删除目录下的kv
	go func() {
		for {
			// 存一下
			kv.Put(context.TODO(), "/cron/job7", "i am job7")

			// 删一下
			kv.Delete(context.TODO(), "/cron/job7")

			// 休息1秒
			time.Sleep(1 * time.Second)
		}
	}()

	// 获取当前/cron/job6的值, 然后监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// 如果Get时刻kv存在, 则记录下来
	if len(getResp.Kvs) != 0 {
		curVal = string(getResp.Kvs[0].Value)
	}

	// 用于监听kv变化的watcher
	watcher = clientv3.NewWatcher(client)

	// 演示10秒, 然后终止watch
	time.AfterFunc(10 * time.Second, func() {
		watcher.Close()
	})

	// 从Get操作时etcd的集群版本号开始监听后续变化
	watchStartRev = getResp.Header.Revision + 1

	fmt.Println("从该版本监听后续变化:", watchStartRev)
	watchChan = watcher.Watch(context.TODO(), "/cron/job7", clientv3.WithRev(watchStartRev))

	// 处理PUT和DELETE事件
	for watchResp = range watchChan {
		for _, event = range watchResp.Events {
			switch (event.Type) {
			case mvccpb.PUT:
				curVal = string(event.Kv.Value)
				fmt.Println("PUT:", curVal, "Revision:", event.Kv.ModRevision)
			case mvccpb.DELETE:
				curVal = ""
				fmt.Println("DEL:", curVal, "Revision:", event.Kv.ModRevision)
			}
		}
	}

	// 监听被Close, 所以watchChan被关闭
	fmt.Println("停止监听")
}

func main() {
	demo7()
}