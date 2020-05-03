package models

//禁止的命令
type Ban struct {
	Id         int
	Code       string//命令
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
	Status     int//状态  0-正常  1-删除
}

func (ban *Ban) TableName() string {
	return TableName("task_ban")
}