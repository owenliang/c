package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
	"time"
)

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

func main() {
	demo2()
}