package controllers

import (
	"fmt"
	"strings"

	"edgea/enums"
	"edgea/models"
	"edgea/utils"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	controllerName string             //當前控制名稱
	actionName     string             //當前action名稱
	curUser        models.BackendUser //當前用戶信息
}

func (c *BaseController) Prepare() {
	//附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	//從Session裡獲取數據 設置用戶信息
	c.adapterUserInfo()
}

// checkLogin判斷用戶是否登錄，未登錄則跳轉至登錄頁面
// 一定要在BaseController.Prepare()後執行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登錄頁面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="
		//登錄成功後返回的址為當前
		returnURL := c.Ctx.Request.URL.Path
		//如果ajax請求則返回相應的錯碼和跳轉的地址
		if c.Ctx.Input.IsAjax() {
			//由於是ajax請求，因此地址是header裡的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(enums.JRCode302, "請登錄", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)
		c.StopRun()
	}
}

// 判斷某 Controller.Action 當前用戶是否有權訪問
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}
	//從session獲取用戶信息
	user := c.GetSession("backenduser")
	//類型斷言
	v, ok := user.(models.BackendUser)
	if ok {
		//如果是超級管理員，則直接通過
		if v.IsSuper == true {
			return true
		}
		//遍歷用戶所負責的資源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
				return true
			}
		}
	}
	return false
}

// checkLogin判斷用戶是否有權訪問某地址，無權則會跳轉到錯誤頁面
//一定要在BaseController.Prepare()後執行
// 會調用checkLogin
// 傳入的參數為忽略權限控制的Action
func (c *BaseController) checkAuthor(ignores ...string) {
	//先判斷是否登錄
	c.checkLogin()
	//如果Action在忽略列表裡，則直接通用
	for _, ignore := range ignores {
		if ignore == c.actionName {
			return
		}
	}
	hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
	if !hasAuthor {
		utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  無權訪問", c.controllerName, c.actionName, c.curUser.Id))
		//如果沒有權限
		if !hasAuthor {
			if c.Ctx.Input.IsAjax() {
				c.jsonResult(enums.JRCode401, "無權訪問", "")
			} else {
				c.pageError("無權訪問")
			}
		}
	}
}

//從session裡取用戶信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("backenduser")
	if a != nil {
		c.curUser = a.(models.BackendUser)
		c.Data["backenduser"] = a
	}
}

//SetBackendUser2Session 獲取用戶信息（包括資源UrlFor）保存至Session
func (c *BaseController) setBackendUser2Session(userId int) error {
	m, err := models.BackendUserOne(userId)
	if err != nil {
		return err
	}
	//獲取這個用戶能獲取到的所有資源列表
	resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("backenduser", *m)
	return nil
}

// 設置模板
// 第一個參數模板，第二個參數為layout
func (c *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "shared/layout_page.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller這個10個字母
		ctrlName := strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
		actionName := strings.ToLower(c.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	c.Layout = layout
	c.TplName = tplName
}
func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

// 重定向 去錯誤頁
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

// 重定向 去登錄頁
func (c *BaseController) pageLogin() {
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
	c.StopRun()
}
