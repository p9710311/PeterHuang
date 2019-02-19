package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 設置BackendUser表名
func (a *Machine) TableName() string {
	return MachineTBName()
}

// BackendUserQueryParam 用於查詢的類
type MachineQueryParam struct {
	BaseQueryParam
	MachineNumberLike string
	MachineNameLike   string
	BrandLike         string
}

// Machine 實體類
type Machine struct {
	Id                     int
	MacAddress             string
	MachineName            string
	MachineNumber          string
	Brand                  string
	Seq                    int
	CollectionIds          []int                     `orm:"-" form:"CollectiothisFormnIds"`
	MachineMoldScheduleRel []*MachineMoldScheduleRel `orm:"reverse(many)"`
	MachineCollectionRel   []*MachineCollectionRel   `orm:"reverse(many)"`
	MachineDashboardARel   []*MachineDashboardARel   `orm:"reverse(many)"`
	Modified               time.Time                 `orm:"auto_now;type(datetime)"`
	CollectionId           int
}

// CollectionPageList 獲取分頁數據
func MachinePageList(params *MachineQueryParam) ([]*Machine, int64) {
	query := orm.NewOrm().QueryTable(MachineTBName())
	data := make([]*Machine, 0)
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
	query = query.Filter("machine_name__istartswith", params.MachineNameLike)
	query = query.Filter("brand__istartswith", params.BrandLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// CollectionDataList 獲取角色列表
func MachineDataList(params *MachineQueryParam) []*Machine {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := MachinePageList(params)
	return data
}

// CollectionBatchDelete 批量刪除
func MachineBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(MachineTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// CollectionOne 獲取單條
func MachineOne(id int) (*Machine, error) {
	o := orm.NewOrm()
	m := Machine{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
