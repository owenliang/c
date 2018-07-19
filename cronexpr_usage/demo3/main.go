package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
	"time"
)

// 同时调度多个cron表达式
type CronJob struct {
	expr *cronexpr.Expression	// cron表达式
	nextTime time.Time // 下次调度时间
}

func demo3() {
	var (
		cronJob *CronJob
		scheduleTable map[string]*CronJob
		expr *cronexpr.Expression
		now time.Time
	)

	// 初始化map
	scheduleTable = make(map[string]*CronJob)

	// 当前时间
	now = time.Now()

	// 设置几个定时任务
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{expr, expr.Next(now)}
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{expr, expr.Next(now)}
	scheduleTable["job2"] = cronJob

	// 启动调度协程
	go func() {
		var (
			now time.Time
			jobName string
			cronJob *CronJob
		)

		// 持续调度
		for {
			// 当前时间
			now = time.Now()

			// 遍历所有cron, 查看哪个过期了
			for jobName, cronJob = range scheduleTable {
				// 任务过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动协程执行任务
					go func() {
						fmt.Println("执行:", jobName)
					}()

					// 计算下次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间", cronJob.nextTime)
				}
			}

			select {
			// 休眠500毫秒再进行下一次判定
			case <- time.NewTimer(500 * time.Millisecond).C:
			}
		}
	}()

	// 我们观察20秒
	time.Sleep(20 * time.Second)
}

func main() {
	demo3()
}