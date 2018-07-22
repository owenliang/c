package worker

import (
	"github.com/owenliang/c/common"
	"time"
	"fmt"
)

// cron任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent	// etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulePlan	// 任务调度表
	jobExecutingTable map[string]*common.JobExecuteInfo// 任务执行表( 保存正在运行的任务)
	jobResultChan chan *common.JobExecuteResult // 任务结果队列
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) (needSchedule bool) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo *common.JobExecuteInfo
		jobExisted bool
		jobExecuting bool
		err error
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:	// 保存任务
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan // 更新执行计划
		needSchedule = true
	case common.JOB_EVENT_DELETE: // 删除任务
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name);
		}
		needSchedule = true
	case common.JOB_EVENT_KILL: // 杀死任务
		if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			jobExecuteInfo.CancelFunc() // 仅仅触发杀死进程, 任务最终结束状态以回调为准
		}
	}
	return
}

// 处理任务结果
func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	var (
		jobLog *common.JobLog
	)

	// 删除执行中状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	// 发送执行日志
	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName: result.ExecuteInfo.Job.Name,
			Command: result.ExecuteInfo.Job.Command,
			Output: string(result.Output),
			PlanTime: result.ExecuteInfo.PlanTime.UnixNano() / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000,
			StartTime: result.StartTime.UnixNano() / 1000,
			EndTime: result.EndTime.UnixNano() / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		G_logSink.Append(jobLog)
	}
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)

	// 任务正在执行, 跳过本次
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		// fmt.Println("正在运行")
		return
	}

	//  构建任务执行信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	// 保存执行信息
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	G_executor.ExecuteJob(jobExecuteInfo)
	fmt.Println("执行任务:", jobExecuteInfo.Job.Name, jobExecuteInfo.RealTime, jobExecuteInfo.PlanTime)
}

// 尝试调度到期的任务
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan *common.JobSchedulePlan
		now time.Time
		nearTime *time.Time
	)

	// 当前没有任务, 调度挂起即可
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	// 当前时间
	now = time.Now()

	// 遍历所有任务计划
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) { 	// 任务到期
			scheduler.TryStartJob(jobPlan)	// 尝试启动任务
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}

		// 统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	// 下次调度等待间隔
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// 定时任务调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent *common.JobEvent
		scheduleTimer *time.Timer
		needSchedule bool
		scheduleAfter time.Duration
		jobResult *common.JobExecuteResult
	)

	// 初始化调度
	scheduleAfter = scheduler.TrySchedule()

	// 调度定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	for {
		// 是否需要重新调度
		needSchedule = false

		select {
		case jobEvent = <- scheduler.jobEventChan:	 // 监听任务变化
			needSchedule = scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C: // 最近的任务到期
			needSchedule = true
		case jobResult= <- scheduler.jobResultChan: // 任务执行结果
			scheduler.handleJobResult(jobResult)
		}

		// 任务计划有变化, 重新调度
		if needSchedule {
			scheduleAfter = scheduler.TrySchedule()
			//  重置调度间隔
			scheduleTimer.Reset(scheduleAfter)
		}
	}
}

func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan: make(chan *common.JobExecuteResult, 1000),
	}

	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// 推送任务执行结果
func (scheduler* Scheduler) PushJobResult(jobResult *common.JobExecuteResult)  {
	scheduler.jobResultChan <- jobResult
}