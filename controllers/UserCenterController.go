package controllers

import (
	"strings"

	"edgea/enums"
	"edgea/models"
	"edgea/utils"

	"github.com/astaxie/beego/orm"
)

type UserCenterController struct {
	BaseController
}

func (c *UserCenterController) Prepare() {
	//先執行
	c.BaseController.Prepare()
	//如果一個Controller的所有Action都需要登錄驗證，則將驗證放到Prepare
	c.checkLogin()
}
func (c *UserCenterController) Profile() {
	Id := c.curUser.Id
	m, err := models.BackendUserOne(Id)
	if m == nil || err != nil {
		c.pageError("數據無效，請刷新後重試")
	}
	c.Data["hasAvatar"] = len(m.Avatar) > 0
	utils.LogDebug(m.Avatar)
	c.Data["m"] = m
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "usercenter/profile_headcssjs.html"
	c.LayoutSections["footerjs"] = "usercenter/profile_footerjs.html"
}
func (c *UserCenterController) BasicInfoSave() {
	Id := c.curUser.Id
	oM, err := models.BackendUserOne(Id)
	if oM == nil || err != nil {
		c.jsonResult(enums.JRCodeFailed, "數據無效，請刷新後重試", "")
	}
	m := models.BackendUser{}
	//獲取form裡的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "獲取數據失敗", m.Id)
	}
	oM.RealName = m.RealName
	oM.Mobile = m.Mobile
	oM.Email = m.Email
	oM.Avatar = c.GetString("ImageUrl")
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeFailed, "編輯失敗", m.Id)
	} else {
		c.setBackendUser2Session(Id)
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}
func (c *UserCenterController) PasswordSave() {
	Id := c.curUser.Id
	oM, err := models.BackendUserOne(Id)
	if oM == nil || err != nil {
		c.pageError("數據無效，請刷新後重試")
	}
	oldPwd := strings.TrimSpace(c.GetString("UserPwd", ""))
	newPwd := strings.TrimSpace(c.GetString("NewUserPwd", ""))
	confirmPwd := strings.TrimSpace(c.GetString("ConfirmPwd", ""))
	md5str := utils.String2md5(oldPwd)
	if oM.UserPwd != md5str {
		c.jsonResult(enums.JRCodeFailed, "原密碼錯誤", "")
	}
	if len(newPwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "請輸入新密碼", "")
	}
	if newPwd != confirmPwd {
		c.jsonResult(enums.JRCodeFailed, "兩次輸入的新密碼不一致", "")
	}
	oM.UserPwd = md5str
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeFailed, "保存失敗", oM.Id)
	} else {
		c.setBackendUser2Session(Id)
		c.jsonResult(enums.JRCodeSucc, "保存成功", oM.Id)
	}
}
func (c *UserCenterController) UploadImage() {
	//這裡type沒有用，只是為了演示傳值
	stype, _ := c.GetInt32("type", 0)
	if stype > 0 {
		f, h, err := c.GetFile("fileImageUrl")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "上傳失敗", "")
		}
		defer f.Close()
		filePath := "static/upload/" + h.Filename
		// 保存位置在 static/upload, 沒有文件夾要先創建
		c.SaveToFile("fileImageUrl", filePath)
		c.jsonResult(enums.JRCodeSucc, "上傳成功", "/"+filePath)
	} else {
		c.jsonResult(enums.JRCodeFailed, "上傳失敗", "")
	}
}
