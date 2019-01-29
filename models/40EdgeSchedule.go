package models

import "github.com/astaxie/beego/orm"

// TableName 設置BackendUser表名
func (a *Schedule) TableName() string {
	return ScheduleTBName()
}

// ScheduleQueryParam 用於查詢的類
type ScheduleQueryParam struct {
	BaseQueryParam
	MoldNumberLike string
}

// Schedule 實體類
type Schedule struct {
	Id                     int
	Seq                    int
	Qty                    string
	MoldNumber             string
	MachineName            stringmachines
	TimeStart              string
	TimeEnd                string
	MachineMoldScheduleRel []*MachineMoldScheduleRel `orm:"reverse(many)"`
	MoldIds                []int                     `orm:"-" form:"MoldIds"`
	MoldId                 int
	MachineIds             []int `orm:"-" form:"MachineIds"`
	MachineId              int
}

// SchedulePageList 獲取分頁數據
func SchedulePageList(params *ScheduleQueryParam) ([]*Schedule, int64) {
	query := orm.NewOrm().QueryTable(ScheduleTBName())
	data := make([]*Schedule, 0)
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
	query = query.Filter("moldnumber__istartswith", params.MoldNumberLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// ScheduleDataList 獲取角色列表
func ScheduleDataList(params *ScheduleQueryParam) []*Schedule {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "desc"
	data, _ := SchedulePageList(params)
	return data
}

// ScheduleBatchDelete 批量刪除
func ScheduleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(ScheduleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// ScheduleOne 獲取單條
func ScheduleOne(id int) (*Schedule, error) {
	o := orm.NewOrm()
	m := Schedule{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
