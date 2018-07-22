package worker

import (
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"github.com/owenliang/c/common"
)

// 分布式任务锁
type JobLock struct {
	// etcd访问方法
	kv clientv3.KV
	lease clientv3.Lease

	isLocked bool // 是否锁住
	jobName string // 任务名
	leaseId clientv3.LeaseID // 租约ID
	cancelFunc context.CancelFunc	// 终止续租
}

// 初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock){
	jobLock = &JobLock{
		kv: kv,
		lease: lease,
		jobName: jobName,
	}
	return
}

// 尝试上锁
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx context.Context
		cancelFunc context.CancelFunc
		leaseId clientv3.LeaseID
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		txn clientv3.Txn
		lockKey string
		txnResp *clientv3.TxnResponse
	)

	// 创建5秒的租约
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	// context用于终止续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	// 租约的ID
	leaseId = leaseGrantResp.ID

	// 开始自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	// 拉起续租应答协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {	// 续租异常或者主动终止
					goto END
				}
			}
		}
		END:
	}()

	// 创建事务
	txn = jobLock.kv.Txn(context.TODO())

	// 锁路径
	lockKey = common.BuildJobLockKey(jobLock.jobName)

	// 抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	// 锁已被占用,
	if !txnResp.Succeeded {
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}

	// 抢锁成功
	jobLock.isLocked = true
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	return

FAIL:
	cancelFunc() // 取消续租
	jobLock.lease.Revoke(context.TODO(), leaseId) // 删除租约
	return
}

// 释放锁
func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		jobLock.isLocked = false
		jobLock.cancelFunc() // 取消续租
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)	// 删除租约
	}
}