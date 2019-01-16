package controllers

import (
	"strings"

	"edgea/enums"
	"edgea/models"
	"edgea/utils"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	//判斷是否登錄
	c.checkLogin()
	c.setTpl()
}
func (c *HomeController) Index2() {
	//判斷是否登錄
	c.checkLogin()
	c.setTpl()
}
func (c *HomeController) Page404() {
	c.setTpl()
}
func (c *HomeController) Error() {
	c.Data["error"] = c.GetString(":error")
	c.setTpl("home/error.html", "shared/layout_pullbox.html")
}
func (c *HomeController) Login() {

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
	c.setTpl("home/login.html", "shared/layout_base.html")
}
func (c *HomeController) DoLogin() {
	username := strings.TrimSpace(c.GetString("UserName"))
	userpwd := strings.TrimSpace(c.GetString("UserPwd"))
	if len(username) == 0 || len(userpwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "用戶名和密碼不正確", "")
	}
	userpwd = utils.String2md5(userpwd)
	user, err := models.BackendUserOneByUserName(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFailed, "用戶被禁用，請聯繫管理員", "")
		}
		//保存用戶信息到session
		c.setBackendUser2Session(user.Id)
		//獲取用戶信息
		c.jsonResult(enums.JRCodeSucc, "登錄成功", "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "用戶名或者密碼錯誤", "")
	}
}
func (c *HomeController) Logout() {
	user := models.BackendUser{}
	c.SetSession("backenduser", user)
	c.pageLogin()
}
