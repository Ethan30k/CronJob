package models

//服务器分组
type ServerGroup struct {
	Id          int
	GroupName   string//组名
	Description string//说明
	Status      int//1-正常   0-删除
	CreateTime  int64//创建时间
	UpdateTime  int64//更新时间
	CreateId    int//创建者id
	UpdateId    int//更新者id
}

func (servergroup *ServerGroup) TableName() string {
	return TableName("task_server_group")
}