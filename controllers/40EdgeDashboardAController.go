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

// 協作測試
type DashboardAController struct {
	BaseController
}

func (c *DashboardAController) Prepare() {
	c.BaseController.Prepare()
	c.checkLogin()

}
func (c *DashboardAController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
}

//Index 角色管理首頁 Back to UX310uQ
func (c *DashboardAController) Index2() {
	//是否顯示更多查詢條件的按鈕
	c.Data["showMoreQuery"] = false
	//將頁面左邊菜單的某項激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "dashboarda/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "dashboarda/index_footerjs.html"
	//頁面裡按鈕權限控制
	c.Data["canEdit"] = c.checkActionAuthor("DashboardAController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("DashboardAController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("DashboardAController", "Edit2")
}

func (c *DashboardAController) DataGrid() {
	var params models.DashboardAQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	data, total := models.DashboardAList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *DashboardAController) DataList() {
	var params = models.DashboardAQueryParam{}
	//獲取數據列表和總數
	data := models.DashboardADataList(&params)
	//定義返回的數據結構
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、編輯角色界面
func (c *DashboardAController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.DashboardA{}
	user := models.DashboardA{Id: Id}
	var err error
	if Id > 0 {
		// o := orm.NewOrm()
		// err := o.Read(&m)
		// if err != nil {
		// 	c.pageError("數據無效，請刷新後重試")
		// }
		m, err = models.DashboardAOne(Id)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
		o := orm.NewOrm()
		o.LoadRelated(m, "MachineDashboardARel")
		o.Read(&user)
	}

	c.Data["m"] = m
	c.setTpl("dashboarda/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "dashboarda/edit_footerjs.html"

	var moldIds []string
	var machineIds []string
	var machinenamesstring string = strconv.Itoa(user.MachineId)
	var moldstring string = strconv.Itoa(user.MoldId)
	machinenames := []string{machinenamesstring}
	moldnameid := []string{moldstring}
	for _, item := range m.MachineDashboardARel {
		moldIds = append(moldIds, strconv.Itoa(item.Mold.Id))
		machineIds = append(machineIds, strconv.Itoa(item.Machine.Id))

	}
	c.Data["molds"] = strings.Join(moldIds, ",")
	c.Data["machines"] = strings.Join(machineIds, ",")
	c.Data["machinenames"] = strings.Join(machinenames, ",")
	c.Data["moldnames"] = strings.Join(moldnameid, ",")

}

//Edit 添加、編輯角色界面
func (c *DashboardAController) Edit2() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.DashboardA{}

	var err error
	if Id > 0 {
		// o := orm.NewOrm()
		// err := o.Read(&m)
		// if err != nil {
		// 	c.pageError("數據無效，請刷新後重試")
		// }
		m, err = models.DashboardAOne(Id)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
		o := orm.NewOrm()
		o.LoadRelated(m, "MachineDashboardARel")

	}

	c.Data["m"] = m
	c.setTpl("dashboarda/edit2.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "dashboarda/edit2_footerjs.html"
	var moldIds []string
	var machineIds []string

	for _, item := range m.MachineDashboardARel {
		moldIds = append(moldIds, strconv.Itoa(item.Mold.Id))
		machineIds = append(machineIds, strconv.Itoa(item.Machine.Id))
	}
	c.Data["molds"] = strings.Join(moldIds, ",")
	c.Data["machines"] = strings.Join(machineIds, ",")

}

//Save 添加、編輯頁面 保存
func (c *DashboardAController) Save() {
	var err error
	m := models.DashboardA{}
	o := orm.NewOrm()
	//獲取form裡的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)
	}
	//刪除已關聯的歷史數據
	if _, err := o.QueryTable(models.MachineDashboardARelTBName()).Filter("dashboarda__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "刪除歷史關係失敗", "")
	}

	if m.Id == 0 {
		if _, err = o.Insert(&m); err != nil {
			// 	c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
			// } else {
			c.jsonResult(enums.JRCodeFailed, "添加失敗", m.Id)
		}

	} else {
		if _, err := models.DashboardAOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", m.Id)
		}
		if _, err = o.Update(&m); err != nil {
			// 	c.jsonResult(enums.JRCodeSucc, "編輯成功", m.Id)
			// } else {
			c.jsonResult(enums.JRCodeFailed, "編輯失敗", m.Id)
		}
	}

	//添加關係
	var relations []models.MachineDashboardARel

	// var tempval int = 0
	// moldId := len(m.MoldIds)
	// machineId := len(m.MachineIds)
	// if moldId > machineId {
	// 	for i := 0; i < moldId; i++ {
	// 		if machineId < i {
	// 			tempval = machineId
	// 		} else {
	// 			tempval = i
	// 		}
	// 		r := models.Mold{Id: i}
	// 		p := models.Machine{Id: tempval}
	// 		relation := models.MachineDashboardARel{DashboardA: &m, Mold: &r, Machine: &p}
	// 		relations = append(relations, relation)
	// 	}
	// } else {
	// 	for i := 0; i < moldId; i++ {
	// 		if moldId < i {
	// 			tempval = moldId
	// 		} else {
	// 			tempval = i
	// 		}
	// 		r := models.Mold{Id: tempval}
	// 		p := models.Machine{Id: i}
	// 		relation := models.MachineDashboardARel{DashboardA: &m, Mold: &r, Machine: &p}
	// 		relations = append(relations, relation)
	// 	}

	// }
	// var tempval int = 0

	for _, moldId := range m.MoldIds {
		r := models.Mold{Id: moldId}
		p := models.Machine{Id: moldId}
		relation := models.MachineDashboardARel{DashboardA: &m, Mold: &r, Machine: &p}
		relations = append(relations, relation)
	}
	// for _, machineId := range m.MachineIds {
	// 	// r := models.Mold{Id: moldId}
	// 	p := models.Machine{Id: machineId}
	// 	relation := models.MachineDashboardARel{DashboardA: &m, Machine: &p}
	// 	relations[machineId] = append(relations, relation)
	// }

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
func (c *DashboardAController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.DashboardABatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功刪除 %d 項", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "刪除失敗", 0)
	}
}

func (c *DashboardAController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.DashboardAOne(Id)
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

// func (c *DashboardAController) UpdateBoard() {
// 	var params = models.DashboardAQueryParam{}
// 	//獲取數據列表和總數
// 	data := models.DashboardAMain()
// 	//定義返回的數據結構
// 	c.jsonResult(enums.JRCodeSucc, "", data)
// }
