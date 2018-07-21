package common

import (
	"encoding/json"
	"strings"
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
func BuildResponse(errno int , msg string, data *json.RawMessage) (resp []byte, err error) {
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data
	// 序列化
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