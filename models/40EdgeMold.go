package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName
func (a *Mold) TableName() string {
	return MoldTBName()
}

// MoldQueryParam 用於查詢的類
type MoldQueryParam struct {
	BaseQueryParam
	NameLike   string
	NumberLike string
}

// Mold 實體類
type Mold struct {
	Id                     int `form:"Id"`
	Seq                    int
	ClientName             string
	Name                   string `form:"Name"`
	PlasticType            string
	ProductModel           string
	Number                 string
	MachineMoldScheduleRel []*MachineMoldScheduleRel `orm:"reverse(many)"`
	MachineDashboardARel   []*MachineDashboardARel   `orm:"reverse(many)"`
	Created                time.Time                 `orm:"auto_now;type(datetime)"`
}

// MoldPageList 獲取分頁數據
func MoldPageList(params *MoldQueryParam) ([]*Mold, int64) {
	query := orm.NewOrm().QueryTable(MoldTBName())
	data := make([]*Mold, 0)
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
	query = query.Filter("name__istartswith", params.NameLike)
	query = query.Filter("number__istartswith", params.NumberLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// MoldDataList 獲取角色列表
func MoldDataList(params *MoldQueryParam) []*Mold {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := MoldPageList(params)
	return data
}

// MoldBatchDelete 批量刪除
func MoldBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(MoldTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// MoldOne 獲取單條
func MoldOne(id int) (*Mold, error) {
	o := orm.NewOrm()
	m := Mold{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
