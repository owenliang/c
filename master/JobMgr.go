package master

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
}

var (
	// 单例
	G_jobMgr *JobMgr
)


/** API **/

// 初始化任务管理器
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
	)

	// 初始化etcd客户端配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,	// 集群列表
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,	// 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// 读写KV的方法集
	kv = clientv3.NewKV(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
	}
	return
}


