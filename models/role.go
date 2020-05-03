package models

//角色
type Role struct {
	Id             int
	RoleName       string//角色名称
	Detail         string//备注
	ServerGroupIds string//服务器分组id
	TaskGroupIds   string//任务分组id
	CreateId       int//创建者id
	UpdateId       int//更新者id
	Status         int//状态  1-正常  0-删除
	CreateTime     int64//创建时间
	UpdateTime     int64//更新时间
}

func (role *Role) TableName() string {
	return TableName("uc_role")
}