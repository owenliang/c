package main

import (
	"os/exec"
	"fmt"
)

// 执行shell命令, 捕获其输出
func demo2() {
	var (
		cmd *exec.Cmd
		output []byte
		err error
	)

	cmd = exec.Command("/bin/bash", "-c", "ls -l .")
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)	// 打印错误原因
	}
	fmt.Println(string(output))	// 打印运行中的标准输出与错误输出
}

func main() {
	demo2()
}