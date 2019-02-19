package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName
func (a *DashboardA) TableName() string {
	return DashboardATBName()
}

// DashboardAQueryParam 用於查詢的類
type DashboardAQueryParam struct {
	BaseQueryParam
	MachineNumberLike string
}

// DashboardA 結構
type DashboardA struct {
	Id             int
	MacAddress     string
	Seq            int
	StatusColor    string //狀態顏色
	Progress       string //剩餘比例
	PlanProduction string //剩餘數量
	PlanTimes      string //剩餘時間
	MachineNumber  string //設備機號
	MachineName    string //設備名稱
	MoldNumber     string //模號
	Cycletime      string //週期時間
	ExcuteTime     string //訂單已執行時間
	ExcuteQty      string //(訂單)已生產模次
	Qty            string //訂單數量
	FixQty         string
	Group          int
	//各種狀態
	ProgressHoliday  string
	ProgressDowntime string
	ProgressIdle     string
	ProgressAbnormal string
	ProgressRunning  string
	//未使用
	U                    string
	A                    string
	P                    string
	Q                    string
	MachineDashboardARel []*MachineDashboardARel `orm:"reverse(many)"`
	MachineIds           []int                   `orm:"-" form"MachineIds"`
	MachineId            int
	MoldIds              []int `orm:"-" form"MoldIds"`
	MoldId               int
	Modified             time.Time `orm:"auto_now;type(datetime)"`
}

// DashboardAList 獲取分頁數據
func DashboardAList(params *DashboardAQueryParam) ([]*DashboardA, int64) {
	query := orm.NewOrm().QueryTable(DashboardATBName()).OrderBy("Seq")
	data := make([]*DashboardA, 0)

	//默認排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Seq":
		sortorder = "Seq"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("machine_number__istartswith", params.MachineNumberLike)

	total, _ := query.Count()
	query.Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// CollectionDataList 獲取角色列表
func DashboardADataList(params *DashboardAQueryParam) []*DashboardA {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := DashboardAList(params)
	return data
}

// CollectionBatchDelete 批量刪除
func DashboardABatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(DashboardATBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// RoleOne 獲取單條
func DashboardAOne(id int) (*DashboardA, error) {
	o := orm.NewOrm()
	m := DashboardA{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
