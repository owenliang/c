package main

import (
	"github.com/gorhill/cronexpr"
	"fmt"
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

func main() {
	demo1()
}