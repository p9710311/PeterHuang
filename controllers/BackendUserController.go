package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"edgea/enums"
	"edgea/models"
	"edgea/utils"

	"github.com/astaxie/beego/orm"
)

type BackendUserController struct {
	BaseController
}

func (c *BackendUserController) Prepare() {
	//先執行
	c.BaseController.Prepare()
	//如果一個Controller的多數Action都需要權限控制，則將驗證放到Prepare
	c.checkAuthor("DataGrid")
	//如果一個Controller的所有Action都需要登錄驗證，則將驗證放到Prepare
	//權限控制裡會進行登錄驗證，因此這裡不用再作登錄驗證
	// c.checkLogin()

}
func (c *BackendUserController) Index() {
	//是否顯示更多查詢條件的按鈕
	c.Data["showMoreQuery"] = true
	//將頁面左邊菜單的某項激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	//頁面模板設置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "backenduser/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "backenduser/index_footerjs.html"
	//頁面裡按鈕權限控制
	c.Data["canEdit"] = c.checkActionAuthor("BackendUserController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("BackendUserController", "Delete")
}
func (c *BackendUserController) DataGrid() {
	//直接反序化獲取json格式的requestbody裡的值（要求配置文件裡 copyrequestbody=true）
	var params models.BackendUserQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//獲取數據列表和總數
	data, total := models.BackendUserPageList(&params)
	//定義返回的數據結構
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// Edit 添加 編輯 頁面
func (c *BackendUserController) Edit() {
	//如果是Post請求，則由Save處理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.BackendUser{}
	var err error
	if Id > 0 {
		m, err = models.BackendUserOne(Id)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
		o := orm.NewOrm()
		o.LoadRelated(m, "RoleBackendUserRel")
	} else {
		//添加用戶時默認狀態為啟用
		m.Status = enums.Enabled
	}
	c.Data["m"] = m
	//獲取關聯的roleId列表
	var roleIds []string
	for _, item := range m.RoleBackendUserRel {
		roleIds = append(roleIds, strconv.Itoa(item.Role.Id))
	}
	c.Data["roles"] = strings.Join(roleIds, ",")
	c.setTpl("backenduser/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/edit_footerjs.html"
}
func (c *BackendUserController) Save() {
	m := models.BackendUser{}
	o := orm.NewOrm()
	var err error
	//獲取form裡的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)
	}
	//刪除已關聯的歷史數據
	if _, err := o.QueryTable(models.RoleBackendUserRelTBName()).Filter("backenduser__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "刪除歷史關係失敗", "")
	}
	if m.Id == 0 {
		//對密碼進行加密
		m.UserPwd = utils.String2md5(m.UserPwd)
		if _, err := o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失敗", m.Id)
		}
	} else {
		if oM, err := models.BackendUserOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", m.Id)
		} else {
			m.UserPwd = strings.TrimSpace(m.UserPwd)
			if len(m.UserPwd) == 0 {
				//如果密碼為空則不修改
				m.UserPwd = oM.UserPwd
			} else {
				m.UserPwd = utils.String2md5(m.UserPwd)
			}
			//本頁面不修改頭像和密碼，直接將值附給新m
			m.Avatar = oM.Avatar
		}
		if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "編輯失敗", m.Id)
		}
	}
	//添加關係
	var relations []models.RoleBackendUserRel
	for _, roleId := range m.RoleIds {
		r := models.Role{Id: roleId}
		relation := models.RoleBackendUserRel{BackendUser: &m, Role: &r}
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
func (c *BackendUserController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.BackendUserTBName())
	if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功刪除 %d 項", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "刪除失敗", 0)
	}
}
