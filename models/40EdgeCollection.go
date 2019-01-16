package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 設置BackendUser表名
func (a *Collection) TableName() string {
	return CollectionTBName()
}

// CollectionQueryParam 用於查詢的類 Peter first Git
type CollectionQueryParam struct {
	BaseQueryParam
	MacAddressLike string
}

// Collection 實體類
type Collection struct {
	Id         int    `form:"Id"`
	MacAddress string `form:"MacAddress"`
	Seq        int
	WiseIp     string
	Isassign   bool
	Created    time.Time `orm:"auto_now;type(datetime)"`
	// CollectionMachineRel []*CollectionMachineRel `orm:"reverse(many)" json:"-"`
}

// CollectionPageList 獲取分頁數據
func CollectionPageList(params *CollectionQueryParam) ([]*Collection, int64) {
	query := orm.NewOrm().QueryTable(CollectionTBName())
	data := make([]*Collection, 0)
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
	// query = query.Filter("macaddress__istartswith", params.MacAddressLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// CollectionDataList 獲取角色列表
func CollectionDataList(params *CollectionQueryParam) []*Collection {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := CollectionPageList(params)
	return data
}

// CollectionBatchDelete 批量刪除
func CollectionBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(CollectionTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// CollectionOne 獲取單條
func CollectionOne(id int) (*Collection, error) {
	o := orm.NewOrm()
	m := Collection{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
