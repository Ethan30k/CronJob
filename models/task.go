package models

// 0：正常  -1：错误  -2:超时
const (
	TASK_SUCESS=0
	TASK_ERROR=-1
	TSAK_TIMEOUT=-2
)

//任务
type Task struct {
	Id            int
	GroupId       int //任务分组id
	ServerId      int//服务器id，当值为0时，在本地机器执行
	TaskName      string//任务名称
	Description   string//任务描述
	CronSpec      string//时间表达式
	Concurrent    int//表示是否允许一个实例，0表示只允许一个实例
	Command       string//命令
	Timeout       int//超时时间
	ExecuteTimes  int//执行时间
	PrevTime      int64//上一次执行时间
	IsNotify      int//是否通知管理员：0-不通知 1-通知
	NotifyType    int//通知类型：0-邮件通知  1-短信通知
	NotifyUserIds string//被通知人的id
	Status        int//状态   -1：删除  0：停用  1：启用  3：不通过
	CreateTime    int64//创建时间
	CreateId      int//创建人的id
	UpdateTime    int64//更新时间
	UpdateId      int//更新者id
}

func (task *Task) TableName() string {
	return TableName("task")
}