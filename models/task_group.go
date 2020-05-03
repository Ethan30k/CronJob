package models


//任务分组
type Group struct {
	Id          int
	GroupName   string//组名
	Description string//描述
	CreateId    int//创建者id
	CreateTime  int64//创建时间
	UpdateId    int//更新者id
	UpdateTime  int64//更新时间
	Status      int//状态   1-正常  0-删除
}

func (group *Group) TableName() string {
	return TableName("task_group")
}