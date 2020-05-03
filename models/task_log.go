package models

//任务日志
type TaskLog struct {
	Id          int
	TaskId      int//任务id创建时间
	Output      string//输出
	Error       string//错误信息
	Status      int//状态     0：正常  -1：错误  -2:超时
	ProcessTime int//任务执行时间
	CreateTime  int64//
}


func (tasklog *TaskLog) TableName() string {
	return TableName("task_log")
}