package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"edgea/enums"
	"edgea/models"

	"github.com/astaxie/beego/orm"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) Prepare() {
	//先執行
	c.BaseController.Prepare()
	//如果一個Controller的少數Action需要權限控制，則將驗證放到需要控制的Action裡
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一個Controller的所有Action都需要登錄驗證，則將驗證放到Prepare
	//這裡註釋了權限控制，因此這裡需要登錄驗證
	c.checkLogin()
}
func (c *ResourceController) Index() {
	//需要權限控制
	c.checkAuthor()
	//將頁面左邊菜單的某項激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//頁面裡按鈕權限控制
	c.Data["canEdit"] = c.checkActionAuthor("ResourceController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("ResourceController", "Delete")
}

//TreeGrid 獲取所有資源的列表
func (c *ResourceController) TreeGrid() {
	tree := models.ResourceTreeGrid()
	//轉換UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//UserMenuTree 獲取用戶有權管理的菜單、區域列表
func (c *ResourceController) UserMenuTree() {
	userid := c.curUser.Id
	//獲取用戶有權管理的菜單列表（包括區域）
	tree := models.ResourceTreeGridByUserId(userid, 1)
	//轉換UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//ParentTreeGrid 獲取可以成為某節點的父節點列表
func (c *ResourceController) ParentTreeGrid() {
	Id, _ := c.GetInt("id", 0)
	tree := models.ResourceTreeGrid4Parent(Id)
	//轉換UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// UrlFor2LinkOne 使用URLFor方法，將資源表裡的UrlFor值轉成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
	// ResourceController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return c.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}
		return c.URLFor(strs[0], values...)
	}
	return ""
}

//UrlFor2Link 使用URLFor方法，批量將資源表裡的UrlFor值轉成LinkUrl
func (c *ResourceController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

//Edit 資源編輯頁面
func (c *ResourceController) Edit() {
	//需要權限控制
	c.checkAuthor()
	//如果是POST請求，則由Save處理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Resource{}
	var err error
	if Id == 0 {
		m.Seq = 100
	} else {
		m, err = models.ResourceOne(Id)
		if err != nil {
			c.pageError("數據無效，請刷新後重試")
		}
	}
	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}
	//獲取可以成為當前節點的父節點的列表
	c.Data["parents"] = models.ResourceTreeGrid4Parent(Id)
	//轉換地址
	m.LinkUrl = c.UrlFor2LinkOne(m.UrlFor)
	c.Data["m"] = m
	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}

	c.setTpl("resource/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/edit_footerjs.html"
}

//Save 資源添加編輯 保存
func (c *ResourceController) Save() {
	var err error
	o := orm.NewOrm()
	parent := &models.Resource{}
	m := models.Resource{}
	parentId, _ := c.GetInt("Parent", 0)
	//獲取form裡的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)
	}
	//獲取父節點
	if parentId > 0 {
		parent, err = models.ResourceOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.jsonResult(enums.JRCodeFailed, "父節點無效", "")
		}
	}
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

// Delete 刪除
func (c *ResourceController) Delete() {
	//需要權限控制
	c.checkAuthor()
	Id, _ := c.GetInt("Id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "選擇的數據無效", 0)
	}
	query := orm.NewOrm().QueryTable(models.ResourceTBName())
	if _, err := query.Filter("id", Id).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("刪除成功"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "刪除失敗", 0)
	}
}

// Select 通用選擇面板
func (c *ResourceController) Select() {
	//獲取調用者的類別 1表示 角色
	desttype, _ := c.GetInt("desttype", 0)
	//獲取調用者的值
	destval, _ := c.GetInt("destval", 0)
	//返回的資源列表
	var selectedIds []string
	o := orm.NewOrm()
	if desttype > 0 && destval > 0 {
		//如果都大於0,則獲取已選擇的值，例如：角色，就是獲取某個角色已關聯的資源列表
		switch desttype {
		case 1:
			{
				role := models.Role{Id: destval}
				o.LoadRelated(&role, "RoleResourceRel")
				for _, item := range role.RoleResourceRel {
					selectedIds = append(selectedIds, strconv.Itoa(item.Resource.Id))
				}
			}
		}
	}
	c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	c.setTpl("resource/select.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/select_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/select_footerjs.html"
}

//CheckUrlFor 填寫UrlFor時進行驗證
func (c *ResourceController) CheckUrlFor() {
	urlfor := c.GetString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.jsonResult(enums.JRCodeSucc, "解析成功", link)
	} else {
		c.jsonResult(enums.JRCodeFailed, "解析失敗", link)
	}
}
func (c *ResourceController) UpdateSeq() {

	Id, _ := c.GetInt("pk", 0)
	oM, err := models.ResourceOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "選擇的數據無效", 0)
	}
	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	if _, err := orm.NewOrm().Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失敗", oM.Id)
	}
}
