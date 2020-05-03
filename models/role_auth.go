package models

type RoleAuth struct {
	AuthId int `orm:"pk"`
	RoleId int64
}

func (roleauth *RoleAuth) TableName() string {
	return TableName("uc_auth_role")
}