package main

import (
	"os/exec"
	"fmt"
)

// 执行shell命令
func demo1() {
	var (
		cmd *exec.Cmd
		err error
	)

	cmd = exec.Command("/bin/bash", "-c", "echo 1; echo 2; sleep 1; exit 5;")
	err = cmd.Run()
	fmt.Println(err)
}

func main ()  {
	demo1()
}