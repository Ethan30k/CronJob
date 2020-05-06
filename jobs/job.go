package jobs

import (
	"CronJob/models"
	"errors"
	"time"
	"github.com/astaxie/beego"
	"fmt"
	"runtime/debug"
	"strings"
	"strconv"
	"CronJob/notify"
	"runtime"
	"os/exec"
	"bytes"
	"golang.org/x/crypto/ssh"

)

type Job struct {
	id int
	logId int64//日志id
	name string //job名
	task *models.Task
	//duration:超时时间
	//返回值一：任务输出
	//返回值二：错误信息
	//返回值三：出现的错误
	//返回值四：是否超时
	runFunc func(duration time.Duration)(string, string, error, bool)
	//任务状态：大于0表示在执行中
	status int
	//同一个任务是否允许并发执行
	Concurrent bool
}


func (j *Job) Run() {
	//该任务不允许并发执行且该任务正在运行
	if !j.Concurrent && j.status > 0 {
		beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束,本次忽略.", j.id))
		return
	}

	//错误恢复
	defer func() {
		if err := recover(); err != nil {
			beego.Error(err, "\n", string(debug.Stack()))
		}
	}()


	beego.Debug(fmt.Sprintf("开始执行任务: %d", j.id))

	//修改状态
	j.status++
	defer func() {
		j.status--
	}()


	//设置默认的超时时间
	timeout := time.Duration(time.Hour*24)
	//任务的超时时间大于0，则说明用户设置了超时时间，故修改超时时间的默认值
	if j.task.Timeout > 0 {
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
	//任务执行的时间
	log.ProcessTime = int(end.Sub(start) / time.Millisecond)
	//日志创建时间
	log.CreateTime = start.Unix()

	//判断任务是否超时
	if isTimeout {
		log.Status = models.TSAK_TIMEOUT//-2
		log.Error = fmt.Sprintf("任务执行超时%d秒\n---------------\n%s\n",
			timeout/time.Second, cmdOut)
	}else if err != nil {
		log.Status = models.TASK_ERROR//-1
		log.Error = err.Error() + ":" + cmdErr
	}

	TextStatus := []string{
		"<font color='red'>超时</font>",
		"<font color='red'>错误</font>",
		"<font color='green'>正常</font>",
	}

	//0：正常  -1：错误  -2:超时
	status := log.Status + 2

	//任务没有执行成功并且出错需要通知管理员
	if log.Status < 0 && j.task.IsNotify == 1 {
		if j.task.NotifyUserIds != "0" && j.task.NotifyUserIds != "" {
			//被发送邮件的人
			toEmail := ""
			//根据管理员id查询管理员
			adminInfo := AllAdminInfo(j.task.NotifyUserIds)
			//遍历管理员信息
			for _, v := range adminInfo {
				//如果管理员邮箱不为空，则拼接管理员邮箱
				if v.Email != "" {
					toEmail += v.Email + ";"
				}
			}
			//  123.com  456.com
			//123.com;456.com;
			//去除toEmail右侧的分号
			toEmail = strings.TrimRight(toEmail, ";")

			//通知类型为邮件并且被通知人的邮箱不为空
			if j.task.NotifyType == 0 && toEmail != "" {
				//创建邮件的主题
				subject := fmt.Sprintf("CronJob定时任务异常：%s", j.task.TaskName)
				//创建邮件的主题内容
				body := fmt.Sprintf(
					`Hello,定时任务出问题了：
					<p style="font-size:16px;">任务执行详情：</p>
					<p style="display:block; padding:10px; background:#efefef;border:1px solid #e4e4e4">
					任务 ID：%d<br/>
					任务名称：%s<br/>
					执行时间：%s<br/>
					执行耗时：%f秒<br/>
					执行状态：%s
					</p>
					<p style="font-size:16px;">任务执行输出</p>
					<p style="display:block; padding:10px; background:#efefef;border:1px solid #e4e4e4">
					%s
					</p>
					<br/>
					<br/>
					<p>-----------------------------------------------------------------<br />
					本邮件由CronJob定时系统自动发出，请勿回复<br />
					如果要取消邮件通知，请登录到系统进行设置<br />
					</p>
					`, j.task.Id,
					j.task.TaskName,
					beego.Date(time.Unix(log.CreateTime, 0), "Y-m-d H:i:s"),
					float64(log.ProcessTime)/1000,
					TextStatus[status],
					log.Error)
				mailtype := "html"//邮件内容的类型

				//将邮件发送到管道
				ok := notify.SendToChan(toEmail, subject, body, mailtype)
				if !ok {
					fmt.Println("发送邮件错误", toEmail)
				}
			}
		}
	}


	//将日志插入到数据库
	models.TaskLogAdd(log)

	j.task.PrevTime = time.Now().Unix()
	j.task.ExecuteTimes++
	j.task.Update("PrevTime", "ExecuteTimes")

}

//  2,3,4
func AllAdminInfo(adminIds string) []*models.Admin {
	//定义切片，用于存储过滤条件
	Filters := make([]interface{}, 0)
	//过滤状态正常的管理员
	Filters = append(Filters, "status", 1)
	var notifyUserIds []int
	//管理员id不是默认值并且不是空字符串
	if adminIds != "0" && adminIds != "" {
		//通过逗号切割adminIds
		notifyUserIdStr := strings.Split(adminIds, ",")
		//遍历切片
		for _, v := range notifyUserIdStr {
			//将字符串id转换为整形id
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
		//追加过滤条件
		Filters = append(Filters, "id__in", notifyUserIds)
	}
	//查询
	Result, _ := models.AdminGetList(1, 1000, Filters)
	return Result
}

//根据task创建job
func NewJobFromTask(task *models.Task) (*Job, error){
	if task.Id < 1 {
		return  nil, fmt.Errorf("ToJob:缺少id!")
	}
	//本地执行
	if task.ServerId == 0 {
		job := NewCommandJob(task.Id, task.TaskName, task.Command)
		job.task = task
		job.Concurrent = task.Concurrent == 1
		return job, nil
	}

	//根据服务器id查询服务器
	server, _ := models.TaskSeverGetById(task.ServerId )
	if server.Type == 0 {
		job := RemoteCommandJobByPassword(task.Id, task.TaskName, task.Command, server)
		job.task = task
		job.Concurrent = task.Concurrent == 1
		return job, nil
	}
	return nil, errors.New("Job创建失败!")
}

//返回值一：任务输出
//返回值二：错误信息
//返回值三：出现的错误
//返回值四：是否超时
func NewCommandJob(id int, name string, command string) *Job {
	job := &Job{
		id:id,
		name:name,
	}
	job.runFunc = func(timeout time.Duration)(string, string, error, bool) {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("CMD", "/C", command)
		}else {
			cmd = exec.Command("sh", "-C", command)
		}
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd.Stdout = bufOut//重定向输出
		cmd.Stderr = bufErr//重定向错误
		//执行命令
		cmd.Start()
		//等待命令结束，释放资源
		err, isTimeout := runCmdWithTimeout(cmd, timeout)
		return bufOut.String(), bufErr.String(), err,isTimeout
	}
	return job
}


//远程执行命令
func RemoteCommandJobByPassword(id int, name string, command string, servers *models.TaskServer) *Job {
	var (
		auth []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		client *ssh.Client
		err error
		session *ssh.Session
	)
	job := &Job{
		id:id,
		name:name,
	}
	job.runFunc = func(timeout time.Duration)(string, string, error, bool) {
		//ssh.AuthMethod里面存放了ssh的认证方式，如果使用密码进行认证，
		//需要用ssh.Password()来加载密码
		auth = make([]ssh.AuthMethod, 0)
		//追加密码
		auth = append(auth, ssh.Password(servers.Password))
		//创建客户端配置
		clientConfig = &ssh.ClientConfig{
			User:servers.ServerAccount,
			Auth: auth,
		}

		//创建client
		if client, err = ssh.Dial("tcp", servers.ServerIp, clientConfig); err != nil {
			return "", "", err, false
		}

		defer client.Close()

		//创建连接
		if session, err =client.NewSession(); err != nil {
			return "", "", err, false
		}

		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		session.Stdout = bufOut//重定向输出
		session.Stderr = bufErr//重定向错误

		//执行定时任务
		if err = session.Run(command); err != nil {
			return "", "", err, false
		}
		return bufOut.String(), bufErr.String(), err, false
	}
	return job
}




































