package common

import "encoding/json"

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