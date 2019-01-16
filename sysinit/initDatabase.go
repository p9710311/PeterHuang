package sysinit

import (
	_ "edgea/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

//初始化數據連接
func InitDatabase() {
	//讀取配置文件，設置數據庫參數
	//數據庫類別
	dbType := beego.AppConfig.String("db_type")
	//連接名稱
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//數據庫名稱
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//數據庫連接用戶名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//數據庫連接用戶名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//數據庫IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//數據庫端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?charset="+dbCharset, 30)
	}
	//如果是開發模式，則顯示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")
	//自動建表
	orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}
}
