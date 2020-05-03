package models

//管理员
type Admin struct {
	Id         int
	LoginName  string//登录名称
	RealName   string//真实姓名
	Password   string//密码
	RoleIds    string//角色id字符串
	Phone      string//联系电话
	Email      string//邮箱
	Salt       string//密码盐
	LastLogin  int64//最后登录时间
	LastIp     string//最后登录的ip
	Status     int//状态 1-正常  0-禁用
	CreateId   int//创建者id
	UpdateId   int//更新者id
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
}


func (admin *Admin) TableName() string {
	return TableName("uc_admin")
}