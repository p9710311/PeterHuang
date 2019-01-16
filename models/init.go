package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(BackendUser), new(Resource), new(Role), new(RoleResourceRel), new(RoleBackendUserRel), new(Schedule), new(DashboardA), new(Collection), new(Machine), new(Mold), new(MachineMoldScheduleRel))
}

// TableName 下面是統一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_prefix")
	return prefix + name
}

// BackendUserTBName 獲取 BackendUser 對應的表名稱
func BackendUserTBName() string {
	return TableName("backend_user")
}

// ResourceTBName 獲取 Resource 對應的表名稱
func ResourceTBName() string {
	return TableName("resource")
}

// RoleTBName 獲取 Role 對應的表名稱
func RoleTBName() string {
	return TableName("role")
}

// RoleResourceRelTBName 角色與資源多對多關係表
func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

// RoleBackendUserRelTBName 角色與用戶多對多關係表
func RoleBackendUserRelTBName() string {
	return TableName("role_backenduser_rel")
}

// DashboardATBName
func DashboardATBName() string {
	return TableName("40_dashboard_a")
}

// MachineTBName
func MachineTBName() string {
	return TableName("40_machine")
}

// CollectionTBName
func CollectionTBName() string {
	return TableName("40_collection")
}

// ScheduleTBName
func ScheduleTBName() string {
	return TableName("40_schedule")
}

// MoldTBName
func MoldTBName() string {
	return TableName("40_mold")
}

// CollectionMachineRelTBName
func CollectionMachineRelTBName() string {
	return TableName("40_collection_machine")
}

// ScheduleRel
func MachineMoldScheduleRelTBName() string {
	return TableName("40_machine_mold_scheule")
}
