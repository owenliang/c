package master

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/owenliang/c/common"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
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
		lease clientv3.Lease
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

	// 用于管理lease租约
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
		lease: lease,
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
	jobKey = common.JOB_SAVE_DIR + job.Name

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

// 删除任务
func (jobMgr *JobMgr)DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey string
		delResp *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	// etcd的任务保存路径
	jobKey = common.JOB_SAVE_DIR + name

	// 删除etcd中的任务
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return // 出错
	}

	// 返回删除前的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 旧值非法, 可以忽略, 返回nil即可
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 列举任务
func (jobMgr *JobMgr)ListJobs() (jobList []*common.Job, err error) {
	var (
		dirKey string
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		job *common.Job
	)

	// etcd的任务保存目录
	dirKey = common.JOB_SAVE_DIR

	// 获取目录下所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 解析所有任务
	jobList = make([]*common.Job, 0)
	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		// 格式异常的任务忽略
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			continue
		}
		jobList = append(jobList, job)
	}

	err = nil
	return
}

// 杀死任务
func (JobMgr *JobMgr)KillJob(name string) (err error) {
	var (
		killerKey string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		// putResp *clientv3.PutResponse
	)

	// etcd通知杀死任务的key
	killerKey = common.JOB_KILLER_DIR + name

	// 任意设置一个过期时间, 只是为了能够过期回收
	if leaseGrantResp, err = JobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	// 租约的ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记位
	if _, err = JobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}

	return
}