package common

import (
	"encoding/json"
	"strings"
	"github.com/gorhill/cronexpr"
	"time"
	"golang.org/x/net/context"
)

// 应答固定协议
type Response struct {
	Errno int `json:"errno"`	// 错误码
	Msg string `json:"msg"`	// 错误原因
	Data *json.RawMessage	`json:"data"`	// 数据
}

// 定时任务
type Job struct {
	Name string `json:"name"`	// 任务名, 全局唯一
	Command string `json:"command"`// shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// 任务变化事件
type JobEvent struct {
	EventType int	// JOB_EVENT_SAVE, JOB_EVENT_DELETE
	Job *Job	// 任务信息
}

// 任务调度计划
type JobSchedulePlan struct {
	Job *Job	// 任务信息
	Expr *cronexpr.Expression	// cron表达式
	NextTime time.Time // 下次调度时间
}

// 任务执行计划
type JobExecuteInfo struct {
	Job *Job // 任务信息
	PlanTime time.Time // 理论调度时间
	RealTime time.Time // 实际调度时间
	CancelCtx context.Context // 任务执行用的context
	CancelFunc context.CancelFunc //  取消任务的cancel函数
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo	// 执行信息
	Output []byte // 脚本输出
	Err error // 脚本失败原因
	StartTime time.Time // 启动时间
	EndTime time.Time // 结束时间
}

// 任务执行日志
type JobLog struct {
	JobName string `bson:"jobName"`// 任务名字
	Command string `bson:"command"`// 脚本命令
	Err string `bson:"err"`// 错误原因
	Output string  `bson:"output"` // shell输出内容
	PlanTime int64 `bson:"planTime"` // 计划开始时间
	ScheduleTime int64 `bson:"scheduleTime"` // 实际调度时间
	StartTime int64 `bson:"startTime"` // 开始执行时间(微秒)
	EndTime int64 `bson:"endTime"` //  结束执行时间
}

// 任务日志查询过滤参数
type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

//  任务日志查询排序
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`	// 按任务开始时间排序
}

// 构造执行计划
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job: jobSchedulePlan.Job,	// 任务信息
		PlanTime: jobSchedulePlan.NextTime, // 计划调度时间
		RealTime: time.Now(),  // 实际调度时间
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

// 构造调度计划
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	// 解析cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	// 生成调度计划
	jobSchedulePlan = &JobSchedulePlan{
		Job: job,
		Expr: expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

// 反序列化任务
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// 构造事件
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job: job,
	}
}

// 构建应答
func BuildResponse(errno int , msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
		bytes []byte
		rawMsg json.RawMessage
	)

	// 序列化data
	if bytes, err = json.Marshal(data); err != nil {
		return
	}
	rawMsg = json.RawMessage(bytes)

	response.Errno = errno
	response.Msg = msg
	response.Data = &rawMsg

	// 序列化整个应答
	resp, err = json.Marshal(response)
	return
}

// 提取任务名
func ExtractJobName(jobKey string) (string) {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 提取要杀死的任务名
func ExtractKillerName(killerKey string) (string) {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// 构建任务锁路径
func BuildJobLockKey(jobName string) (string){
	return JOB_LOCK_DIR + jobName
}