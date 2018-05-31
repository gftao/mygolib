package main

import (
	"fmt"
	"golib/modules/logr"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 5 {
		fmt.Println("usage sleep runMode logFileConfig  可执行程序名")
		return
	}

	//death := make(chan os.Signal, 1)
	//signal.Notify(death, os.Kill, os.Interrupt, syscall.SIGABRT, syscall.SIGQUIT )
	sleepTime, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("输入参数错误:", os.Args[1])
		return
	}
	runMode := os.Args[2]
	logFileConf := os.Args[3]
	execName := os.Args[4]
	err = logr.InitByLogFile(runMode, logFileConf)
	if err != nil {
		fmt.Println("初始化日志失败", err)
		return
	}
	//execName := "ls"

	stdOut := logr.LogBuffer{Name:"stdOut"}
	stdErr := logr.LogBuffer{Name:"stdErr"}

	for {
		var cmd *exec.Cmd
		if len(os.Args) > 5 {
			cmd = exec.Command(execName, os.Args[5:]...)
		} else {
			cmd = exec.Command(execName)
		}
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err := cmd.Start()
		if err != nil {
			logr.Error("启动程序出错", err)
			return
		}

		err = cmd.Wait()
		if err != nil {
			logr.Error("程序执行出错", err)
		} else {
			logr.Error("程序执行完成")
		}
		logr.Error("程序err:", err)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
