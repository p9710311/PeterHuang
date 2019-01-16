package models

import (
	"time"
)

// RoleBackendUserRel 角色與用戶關係
type MachineMoldScheduleRel struct {
	Id int
	// Machine  *Machine  `orm:"rel(fk)"`  //外鍵
	Mold     *Mold     `orm:"rel(fk)" ` // 外鍵
	Schedule *Schedule `orm:"rel(fk)" `
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}

// TableName 設置表名
func (a *MachineMoldScheduleRel) TableName() string {
	return MachineMoldScheduleRelTBName()
}
