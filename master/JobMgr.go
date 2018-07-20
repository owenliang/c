package master

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/owenliang/c/common"
	"encoding/json"
	"golang.org/x/net/context"
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

// 保存任务
func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job
	)

	// etcd的任务保存路径
	jobKey = "/cron/jobs/" + job.Name

	// 任务信息序列化
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	// 保存任务到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	// 如果是覆盖更新, 那么返回旧值
	if putResp.PrevKv != nil {
		// 旧值非法, 可以忽略, 返回nil即可
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			return
		}
		oldJob = &oldJobObj
	}

	return
}

