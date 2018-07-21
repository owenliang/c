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
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) (err error) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo *common.JobExecuteInfo
		jobExisted bool
		jobExecuting bool
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:	// 保存任务
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan // 更新执行计划
	case common.JOB_EVENT_DELETE: // 删除任务
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name);
		}
	case common.JOB_EVENT_KILL: // 杀死任务
		if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			// TODO: 任务正在执行, 杀死它
			jobExecuteInfo = jobExecuteInfo
		}
	}
	return
}

// 重新调度任务
func (scheduler *Scheduler) ReScheduleJobPlan(jobPlan *common.JobSchedulePlan) {
	jobPlan.NextTime = jobPlan.Expr.Next(time.Now())
}

// 定时任务调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent *common.JobEvent
		scheduleTimer *time.Timer
		lastScheduleTime time.Time	// 上次调度时间
		now time.Time
		scheduleIdle time.Duration
		jobPlan *common.JobSchedulePlan
		jobName string
	)

	// 每间隔100毫秒调度一次
	lastScheduleTime = time.Now()
	// TODO: 根据最近一个任务的时间, 计算下次调度时间
	scheduleTimer = time.NewTimer(common.SCHEDULE_PERIOD * time.Millisecond)

	for {
		select {
		case jobEvent = <- scheduler.jobEventChan:	 // 监听任务变化
			scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C: // 等待下个调度周期
		}

		//  判断是否到达调度周期
		now = time.Now()
		scheduleIdle = now.Sub(lastScheduleTime)
		if scheduleIdle < common.SCHEDULE_PERIOD * time.Millisecond {  // 还没到期, 继续等待
			scheduleTimer.Reset(common.SCHEDULE_PERIOD - scheduleIdle)
			continue
		}

		// TODO: 调度任务
		for jobName, jobPlan = range scheduler.jobPlanTable {
			// 任务到期
			if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
				fmt.Println("执行任务:", jobName)
				scheduler.ReScheduleJobPlan(jobPlan)
			}
		}
		lastScheduleTime = time.Now()
		scheduleTimer.Reset(common.SCHEDULE_PERIOD  * time.Millisecond)
	}
}

func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
	}

	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}