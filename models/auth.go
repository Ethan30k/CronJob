package models

//权限
type Auth struct {
	Id         int//主键
	Pid        int//父级id
	AuthName   string//权限名称
	AuthUrl    string//url地址
	Sort       int//排序
	Icon       string
	IsShow     int//是否隐藏
	UserId     int//操作者id
	CreateId   int//创建者id
	UpdateId   int//更新者id
	Status     int//状态 1-正常    0-删除
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
}

func (auth *Auth) TableName() string {
	return TableName("uc_auth")
}

