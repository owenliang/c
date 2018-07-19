package main

import (
	"os/exec"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

// 强制杀死shell命令
type result struct { // shell执行的错误与输出
	err error
	output []byte
}

// 在协程中
func demo3() {
	var (
		ctx context.Context
		cancelFunc context.CancelFunc
		finishedJobs chan *result = make(chan *result, 100)
		res *result
	)

	// 用于通知杀死shell
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 在协程中执行shell
	go func() {
		var (
			cmd *exec.Cmd
			output []byte
			err error
		)

		// 模拟shell脚本执行5秒
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 5")

		// 捕获脚本输出
		output, err = cmd.CombinedOutput()

		// 投递脚本执行结果
		finishedJobs <- &result{err, output}
	}()

	// 1秒后杀死任务
	time.Sleep(1 * time.Second)
	cancelFunc()

	// 等待任务退出
	res = <- finishedJobs

	// 打印任务结果
	fmt.Println(res.err, string(res.output))
}

func main() {
	demo3()
}