package common

const (
	// 任务保存目录
	JOB_SAVE_DIR = "/cron/jobs/"

	// 任务强杀标记
	JOB_KILLER_DIR = "/cron/killer/"

	// 保存任务事件
	JOB_EVENT_SAVE = 1

	// 删除任务事件
	JOB_EVENT_DELETE = 2

	// 杀死任务事件
	JOB_EVENT_KILL = 3

	// 调度周期(毫秒)
	SCHEDULE_PERIOD = 100
)
