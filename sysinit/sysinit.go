package sysinit

import (
	"edgea/utils"

	"github.com/astaxie/beego"
)

func init() {
	//啟用Session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//初始化日誌
	utils.InitLogs()
	//初始化緩存
	utils.InitCache()
	//初始化數據庫
	InitDatabase()
}
