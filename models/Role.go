package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 設置表名
func (a *Role) TableName() string {
	return RoleTBName()
}

// RoleQueryParam 用於搜索的類
type RoleQueryParam struct {
	BaseQueryParam
	NameLike string
}

// Role 用戶角色 實體類
type Role struct {
	Id                 int    `form:"Id"`
	Name               string `form:"Name"`
	Seq                int
	RoleResourceRel    []*RoleResourceRel    `orm:"reverse(many)" json:"-"` // 設置一對多的反向關係
	RoleBackendUserRel []*RoleBackendUserRel `orm:"reverse(many)" json:"-"` // 設置一對多的反向關係
}

// RolePageList 獲取分頁數據
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	data := make([]*Role, 0)
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
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// RoleDataList 獲取角色列表
func RoleDataList(params *RoleQueryParam) []*Role {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := RolePageList(params)
	return data
}

// RoleBatchDelete 批量刪除
func RoleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// RoleOne 獲取單條
func RoleOne(id int) (*Role, error) {
	o := orm.NewOrm()
	m := Role{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
