package models

import (
	"time"
)

// RoleBackendUserRel 角色與用戶關係
type RoleBackendUserRel struct {
	Id          int
	Role        *Role        `orm:"rel(fk)"`  //外鍵
	BackendUser *BackendUser `orm:"rel(fk)" ` // 外鍵
	Created     time.Time    `orm:"auto_now_add;type(datetime)"`
}

// TableName 設置表名
func (a *RoleBackendUserRel) TableName() string {
	return RoleBackendUserRelTBName()
}
