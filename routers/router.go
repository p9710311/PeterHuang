package routers

import (
	"edgea/controllers"

	"github.com/astaxie/beego"
)

func init() {

	// 用戶角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	// 資源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	// 後臺用戶路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "POST:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get,Post:Edit")
	beego.Router("/backenduser/delete", &controllers.BackendUserController{}, "Post:Delete")

	// 後臺用戶中心
	beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/usercenter/passwordsave", &controllers.UserCenterController{}, "Post:PasswordSave")

	// Home
	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")
	beego.Router("/", &controllers.HomeController{}, "*:Index")

	// DashboardA
	beego.Router("/dashboarda/index", &controllers.DashboardAController{}, "*:Index")
	beego.Router("/dashboarda/index2", &controllers.DashboardAController{}, "*:Index2")
	beego.Router("/dashboarda/datagrid", &controllers.DashboardAController{}, "Get,Post:DataGrid")
	beego.Router("/dashboarda/edit/?:id", &controllers.DashboardAController{}, "Get,Post:Edit")
	beego.Router("/dashboarda/edit2/?:id", &controllers.DashboardAController{}, "Get,Post:Edit2")
	beego.Router("/dashboarda/delete", &controllers.DashboardAController{}, "Post:Delete")
	beego.Router("/dashboarda/datalist", &controllers.DashboardAController{}, "Post:DataList")
	// beego.Router("/dashboarda/allocate", &controllers.DashboardAController{}, "Post:Allocate")
	beego.Router("/dashboarda/updateseq", &controllers.DashboardAController{}, "Post:UpdateSeq")

	// Collection
	beego.Router("/collection/index", &controllers.CollectionController{}, "*:Index")
	beego.Router("/collection/datagrid", &controllers.CollectionController{}, "Get,Post:DataGrid")
	beego.Router("/collection/edit/?:id", &controllers.CollectionController{}, "Get,Post:Edit")
	beego.Router("/collection/delete", &controllers.CollectionController{}, "Post:Delete")
	beego.Router("/collection/datalist", &controllers.CollectionController{}, "Post:DataList")
	beego.Router("/collection/updateseq", &controllers.CollectionController{}, "Post:UpdateSeq")

	// Machine
	beego.Router("/machine/index", &controllers.MachineController{}, "*:Index")
	beego.Router("/machine/datagrid", &controllers.MachineController{}, "Get,Post:DataGrid")
	beego.Router("/machine/edit/?:id", &controllers.MachineController{}, "Get,Post:Edit")
	beego.Router("/machine/edit2/?:id", &controllers.MachineController{}, "Get,Post:Edit2")
	beego.Router("/machine/delete", &controllers.MachineController{}, "Post:Delete")
	beego.Router("/machine/datalist", &controllers.MachineController{}, "Post:DataList")
	beego.Router("/machine/updateseq", &controllers.MachineController{}, "Post:UpdateSeq")

	// Mold
	beego.Router("/mold/index", &controllers.MoldController{}, "*:Index")
	beego.Router("/mold/datagrid", &controllers.MoldController{}, "Get,Post:DataGrid")
	beego.Router("/mold/edit/?:id", &controllers.MoldController{}, "Get,Post:Edit")
	beego.Router("/mold/delete", &controllers.MoldController{}, "Post:Delete")
	beego.Router("/mold/datalist", &controllers.MoldController{}, "Post:DataList")
	beego.Router("/mold/updateseq", &controllers.MoldController{}, "Post:UpdateSeq")

	// Schedule
	beego.Router("/schedule/index", &controllers.ScheduleController{}, "*:Index")
	beego.Router("/schedule/datagrid", &controllers.ScheduleController{}, "Get,Post:DataGrid")
	beego.Router("/schedule/edit/?:id", &controllers.ScheduleController{}, "Get,Post:Edit")
	beego.Router("/schedule/edit2/?:id", &controllers.ScheduleController{}, "Get,Post:Edit2")
	beego.Router("/schedule/delete", &controllers.ScheduleController{}, "Post:Delete")
	beego.Router("/schedule/datalist", &controllers.ScheduleController{}, "Post:DataList")
	beego.Router("/schedule/updateseq", &controllers.ScheduleController{}, "Post:UpdateSeq")
}
