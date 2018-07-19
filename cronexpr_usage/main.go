package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
	"time"
)

/*
Field name     Mandatory?   Allowed values    Allowed special characters
----------     ----------   --------------    --------------------------
Seconds        No           0-59              * / , -
Minutes        Yes          0-59              * / , -
Hours          Yes          0-23              * / , -
Day of month   Yes          1-31              * / , - L W
Month          Yes          1-12 or JAN-DEC   * / , -
Day of week    Yes          0-6 or SUN-SAT    * / , - L #
Year           No           1970–2099         * / , -
*/

// cron表达式的解析
func demo1() {
	var (
		expr *cronexpr.Expression
		err error
	)

	// 哪分钟(0-59), 哪小时(0-23), 月内哪天(1-31), 哪月(1-12), 周内哪天(0-6)

	// 0点5分执行1次
	if expr, err = cronexpr.Parse("5 0 * * *"); err != nil {
		fmt.Println(err)
	}

	// 每隔5分钟执行1次
	if expr, err = cronexpr.Parse("*/5 * * * *"); err != nil {
		fmt.Println(err)
	}

	expr = expr
}

// 计算下次执行时间
func demo2() {
	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)

	// 当前时间
	now = time.Now()

	// 每小时的第5分钟执行
	if expr, err = cronexpr.Parse("5 * * * *"); err != nil {
		fmt.Println(err)
	}

	// 计算cron表达式的下次触发时间
	nextTime = expr.Next(now)
	fmt.Println(nextTime)
}

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
	demo1()

	demo2()

	demo3()
}