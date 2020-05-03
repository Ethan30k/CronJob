package jobs

import (
	"CronJob/models"
	"fmt"
	"github.com/astaxie/beego"
	"runtime/debug"
	"time"
)

type Job struct {
	id int
	logId int64		//日志id
	name string		//job名
	task *models.Task
	//duration:超时时间
	//返回值一：任务输出
	//返回值二：错误信息
	//返回值三：出现的错误
	//返回值四：是否超时
	runFunc func(duration time.Duration)(string,string,error,bool)
	//任务状态：大于0表示在执行中
	status int
	//同一个任务是否允许并发执行
	Concurrent bool
}

func (j *Job) Run()  {
	//该任务不允许并发执行且该任务正在运行
	if !j.Concurrent && j.status >0{
		beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次忽略", j.id))
		return
	}

	//错误恢复
	defer func() {
		if err := recover();err!=nil{
			beego.Error(err, "\n",string(debug.Stack()))
		}
	}()

	beego.Debug(fmt.Sprintf("开始执行任务：%d", j.id))
	//修改状态
	j.status++
	defer func() {
		j.status--
	}()

	//设置默认的超时时间
	timeout := time.Duration(24 * time.Hour)
	//任务的超时时间大于0，则说明用户设置了超时时间，故修改超时时间的默认值
	if j.task.Timeout >0{
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	//任务开始执行时间
	start := time.Now()
	//执行定时任务
	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)

	//创建日志对象
	log := new(models.TaskLog)
	log.TaskId = j.id
	log.Output = cmdOut
	log.Error = cmdErr
	//任务执行结束时间
	end := time.Now()
	//任务执行时间
	log.ProcessTime = int(end.Sub(start)/time.Millisecond)
	//日志创建时间
	log.CreateTime = start.Unix()

	//判断任务是否超时
	if isTimeout{
		log.Status = models.TSAK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超时%d秒\n---------------\n%s\n",
			timeout/time.Second)
	}else if err != nil{
		log.Status = models.TASK_ERROR
		log.Error = err.Error()+":" +cmdErr
	}

}