package jobs

import (
	"CronJob/models"
	"os/exec"
	"time"
	"github.com/astaxie/beego"
	"fmt"

)

func InitJob() {
	//查询所有已经启用的任务
	list,_ := models.TaskGetList(1, 10000000, "status", 1)
	//遍历任务列表
	for _, task := range list {
		job, err := NewJobFromTask(task)
		if err != nil {
			beego.Error("InitJobs:", err.Error())
			continue
		}
		AddJob(task.CronSpec, job)
	}

	fmt.Println("Job初始化了")
	fmt.Println("-------------------------------------------------")
	fmt.Println(len(list))
	fmt.Println("-------------------------------------------------")
}

//error：错误信息
//bool:是否超时
func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		//等待命令结束，释放资源
		done <- cmd.Wait()
	}()
	var err error
	select {
	//超时处理
	case <- time.After(timeout):
		beego.Warn("任务执行超时，进程将被强制杀死:%d", cmd.Process.Pid)
		go func() {
			<- done
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("进程无法杀掉：%d, 错误信息: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <- done:
		return err, false
	}

}