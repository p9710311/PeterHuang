package models

import "time"

// RoleResourceRel 角色與資源關係表
type RoleResourceRel struct {
	Id       int
	Role     *Role     `orm:"rel(fk)"`  //外鍵
	Resource *Resource `orm:"rel(fk)" ` // 外鍵
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}

// TableName 設置表名
func (a *RoleResourceRel) TableName() string {
	return RoleResourceRelTBName()
}
