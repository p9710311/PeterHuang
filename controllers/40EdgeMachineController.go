package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"edgea/enums"
	"edgea/models"

	"github.com/astaxie/beego/orm"
)

//MachineController 角色管理
type MachineController struct {
	BaseController
}

//Prepare 參考beego官方文檔說明
func (c *MachineController) Prepare() {
	c.BaseController.Prepare()
	c.checkLogin()
}

//Index 角色管理首頁
func (c *MachineController) Index() {
	c.Data["showMoreQuery"] = false
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "machine/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "machine/index_footerjs.html"
	c.Data["canEdit"] = c.checkActionAuthor("MachineController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("MachineController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("MachineController", "Allocate")
}

// DataGrid 角色管理首頁 表格獲取數據
func (c *MachineController) DataGrid() {
	//直接反序化獲取json格式的requestbody裡的值
	var params models.MachineQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//獲取數據列表和總數
	data, total := models.MachinePageList(&params)
	//定義返回的數據結構
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *MachineController) DataList() {
	var params = models.MachineQueryParam{}
	//獲取數據列表和總數
	data := models.MachineDataList(&params)
	//定義返回的數據結構
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、編輯角色界面
func (c *MachineController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Machine{}
	user := models.Machine{Id: Id}
	var err error
	if Id > 0 {
		m, err = models.MachineOne(Id)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
		o := orm.NewOrm()
		o.LoadRelated(m, "MachineCollectionRel")
		o.Read(&user)
	}
	c.Data["m"] = m
	c.setTpl("machine/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "machine/edit_footerjs.html"
	//獲取關聯
	var collectionIds []string
	// var machIds []string
	var macamesstring string = strconv.Itoa(user.CollectionId)
	macnames := []string{macamesstring}

	for _, item := range m.MachineCollectionRel {
		collectionIds = append(collectionIds, strconv.Itoa(item.Collection.Id))
		// macIds = append(macIds, strconv.Itoa(item.Machine.Id))
	}
	c.Data["collections"] = strings.Join(collectionIds, ",")
	c.Data["macids"] = strings.Join(macnames, ",")

}

//Edit2 添加、編輯角色界面
func (c *MachineController) Edit2() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Machine{Id: Id}

	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
	}
	c.Data["m"] = m
	c.setTpl("machine/edit.new.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "machine/edit_footerjs.new.html"

}

//Save 添加、編輯頁面 保存
func (c *MachineController) Save() {
	var err error
	m := models.Machine{}
	o := orm.NewOrm()
	//獲取form裡的值
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	fmt.Println(c)
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)

	}
	fmt.Println(m)
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)

	}
	fmt.Println(m)

	//刪除已關聯的歷史數據
	if _, err := o.QueryTable(models.MachineCollectionRelTBName()).Filter("machine__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "刪除歷史關係失敗", "")
	}

	if m.Id == 0 {
		if _, err = o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失敗", m.Id)
		}
	} else {
		if _, err := models.MachineOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", m.Id)
		}
		if _, err = o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "編輯失敗", m.Id)
		}
	}

	//添加關係
	var relations []models.MachineCollectionRel
	//var relations2 []models.MachineMoldScheduleRel
	for _, collectionId := range m.CollectionIds {
		r := models.Collection{Id: collectionId}
		relation := models.MachineCollectionRel{Machine: &m, Collection: &r}
		relations = append(relations, relation)
	}

	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失敗", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}

}

//Delete 批量刪除
func (c *MachineController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.MachineBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功刪除 %d 項", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "刪除失敗", 0)
	}
}

//Allocate 給角色分配資源界面
// func (c *MachineController) Allocate() {
// 	collectionId, _ := c.GetInt("id", 0)
// 	strs := c.GetString("ids")

// 	o := orm.NewOrm()
// 	m := models.Collection{Id: collectionId}
// 	if err := o.Read(&m); err != nil {
// 		c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", "")
// 	}
// 	//刪除已關聯的歷史數據
// 	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
// 		c.jsonResult(enums.JRCodeFailed, "刪除歷史關係失敗", "")
// 	}
// 	var relations []models.RoleResourceRel
// 	for _, str := range strings.Split(strs, ",") {
// 		if id, err := strconv.Atoi(str); err == nil {
// 			r := models.Resource{Id: id}
// 			relation := models.RoleResourceRel{Role: &m, Resource: &r}
// 			relations = append(relations, relation)
// 		}
// 	}
// 	if len(relations) > 0 {
// 		//批量添加
// 		if _, err := o.InsertMulti(len(relations), relations); err == nil {
// 			c.jsonResult(enums.JRCodeSucc, "保存成功", "")
// 		}
// 	}
// 	c.jsonResult(0, "保存失敗", "")
// }

func (c *MachineController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.MachineOne(Id)
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
