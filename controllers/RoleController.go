package controllers

import (
	"encoding/json"

	"edgea/enums"
	"edgea/models"

	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

//RoleController 角色管理
type RoleController struct {
	BaseController
}

//Prepare 參考beego官方文檔說明
func (c *RoleController) Prepare() {
	//先執行
	c.BaseController.Prepare()
	//如果一個Controller的多數Action都需要權限控制，則將驗證放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
	//如果一個Controller的所有Action都需要登錄驗證，則將驗證放到Prepare
	//權限控制裡會進行登錄驗證，因此這裡不用再作登錄驗證
	//c.checkLogin()
}

//Index 角色管理首頁
func (c *RoleController) Index() {
	//是否顯示更多查詢條件的按鈕
	c.Data["showMoreQuery"] = false
	//將頁面左邊菜單的某項激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "role/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"
	//頁面裡按鈕權限控制
	c.Data["canEdit"] = c.checkActionAuthor("RoleController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("RoleController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("RoleController", "Allocate")
}

// DataGrid 角色管理首頁 表格獲取數據
func (c *RoleController) DataGrid() {
	//直接反序化獲取json格式的requestbody裡的值
	var params models.RoleQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//獲取數據列表和總數
	data, total := models.RolePageList(&params)
	//定義返回的數據結構
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *RoleController) DataList() {
	var params = models.RoleQueryParam{}
	//獲取數據列表和總數
	data := models.RoleDataList(&params)
	//定義返回的數據結構
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、編輯角色界面
func (c *RoleController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Role{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
	}
	c.Data["m"] = m
	c.setTpl("role/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/edit_footerjs.html"
}

//Save 添加、編輯頁面 保存
func (c *RoleController) Save() {
	var err error
	m := models.Role{}
	//獲取form裡的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失敗", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "編輯成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "編輯失敗", m.Id)
		}
	}

}

//Delete 批量刪除
func (c *RoleController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.RoleBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功刪除 %d 項", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "刪除失敗", 0)
	}
}

//Allocate 給角色分配資源界面
func (c *RoleController) Allocate() {
	roleId, _ := c.GetInt("id", 0)
	strs := c.GetString("ids")

	o := orm.NewOrm()
	m := models.Role{Id: roleId}
	if err := o.Read(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", "")
	}
	//刪除已關聯的歷史數據
	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "刪除歷史關係失敗", "")
	}
	var relations []models.RoleResourceRel
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			r := models.Resource{Id: id}
			relation := models.RoleResourceRel{Role: &m, Resource: &r}
			relations = append(relations, relation)
		}
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", "")
		}
	}
	c.jsonResult(0, "保存失敗", "")
}
func (c *RoleController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.RoleOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "選擇的數據無效", 0)
	}
	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失敗", oM.Id)
	}
}
