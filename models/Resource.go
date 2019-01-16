package models

import (
	"fmt"

	"edgea/utils"

	"github.com/astaxie/beego/orm"
)

// TableName 設置表名
func (a *Resource) TableName() string {
	return ResourceTBName()
}

// Resource 權限控制資源表
type Resource struct {
	Id              int
	Name            string    `orm:"size(64)"`
	Parent          *Resource `orm:"null;rel(fk)"` // RelForeignKey relation
	Rtype           int
	Seq             int
	Sons            []*Resource        `orm:"reverse(many)"` // fk 的反向關係
	SonNum          int                `orm:"-"`
	Icon            string             `orm:"size(32)"`
	LinkUrl         string             `orm:"-"`
	UrlFor          string             `orm:"size(256)" Json:"-"`
	HtmlDisabled    int                `orm:"-"`             //在html裡應用時是否可用
	Level           int                `orm:"-"`             //第幾級，從0開始
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)"` // 設置一對多的反向關係
}

// ResourceOne 獲取單條
func ResourceOne(id int) (*Resource, error) {
	o := orm.NewOrm()
	m := Resource{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ResourceTreeGrid 獲取treegrid順序的列表
func ResourceTreeGrid() []*Resource {
	o := orm.NewOrm()
	query := o.QueryTable(ResourceTBName()).OrderBy("seq", "id")
	list := make([]*Resource, 0)
	query.All(&list)
	return resourceList2TreeGrid(list)
}

// ResourceTreeGrid4Parent 獲取可以成為某個節點父節點的列表
func ResourceTreeGrid4Parent(id int) []*Resource {
	tree := ResourceTreeGrid()
	if id == 0 {
		return tree
	}
	var index = -1
	//找出當前節點所在索引
	for i, _ := range tree {
		if tree[i].Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return tree
	} else {
		tree[index].HtmlDisabled = 1
		for _, item := range tree[index+1:] {
			if item.Level > tree[index].Level {
				item.HtmlDisabled = 1
			} else {
				break
			}
		}
	}
	return tree
}

// ResourceTreeGridByUserId 根據用戶獲取有權管理的資源列表，並整理成teegrid格式
func ResourceTreeGridByUserId(backuserid, maxrtype int) []*Resource {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", backuserid, maxrtype)
	var list []*Resource
	if err := utils.GetCache(cachekey, &list); err == nil {
		return list
	}
	o := orm.NewOrm()
	user, err := BackendUserOne(backuserid)
	if err != nil || user == nil {
		return list
	}
	var sql string
	if user.IsSuper == true {
		//如果是管理員，則查出所有的
		sql = fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, ResourceTBName())
		o.Raw(sql, maxrtype).QueryRows(&list)
	} else {
		//聯查多張表，找出某用戶有權管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T0.resource_id,T2.id,T2.name,T2.parent_id,T2.rtype,T2.icon,T2.seq,T2.url_for
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.resource_id
		WHERE T1.backend_user_id = ? and T2.rtype <= ?  Order By T2.seq asc,T2.id asc`, RoleResourceRelTBName(), RoleBackendUserRelTBName(), ResourceTBName())
		o.Raw(sql, backuserid, maxrtype).QueryRows(&list)
	}
	result := resourceList2TreeGrid(list)
	utils.SetCache(cachekey, result, 30)
	return result
}

// resourceList2TreeGrid 將資源列表轉成treegrid格式
func resourceList2TreeGrid(list []*Resource) []*Resource {
	result := make([]*Resource, 0)
	for _, item := range list {
		if item.Parent == nil || item.Parent.Id == 0 {
			item.Level = 0
			result = append(result, item)
			result = resourceAddSons(item, list, result)
		}
	}
	return result
}

//resourceAddSons 添加子菜單
func resourceAddSons(cur *Resource, list, result []*Resource) []*Resource {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			cur.SonNum++
			item.Level = cur.Level + 1
			result = append(result, item)
			result = resourceAddSons(item, list, result)
		}
	}
	return result
}
