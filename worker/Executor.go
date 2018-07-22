package worker

import (
	"github.com/owenliang/c/common"
	"os/exec"
	"time"
)

// 任务执行器
type Executor struct {

}

var (
	G_executor *Executor
)

// 初始化
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}

// 执行任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	// 启动协程执行shell任务
	go func() {
		var (
			cmd *exec.Cmd
			output []byte
			err error
			startTime time.Time
			endTime time.Time
			result *common.JobExecuteResult
			jobLock *JobLock
		)

		// 执行结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output: make([]byte, 0),
		}

		// 创建锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		// 上锁
		result.StartTime = time.Now()

		if err = jobLock.TryLock(); err != nil {	// 上锁失败
			result.EndTime = time.Now()
			result.Err = err
		} else {
			// 开始时间
			startTime = time.Now()

			// 创建Shell任务
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

			// 执行命令, 捕获输出
			output, err = cmd.CombinedOutput()

			// 结束时间
			endTime = time.Now()

			// 构造结果
			result.StartTime = startTime
			result.EndTime = endTime
			result.Err = err
			result.Output = output

			// 释放锁
			jobLock.Unlock()

			// 在命令行看效果:
			// ETCDCTL_API=3 ./etcdctl watch / --prefix
		}

		// 通知调度器
		G_scheduler.PushJobResult(result)
	}()
}

