package controllers

type MoldAnalyzeController struct {
	BaseController
}

func (c *MoldAnalyzeController) Prepare() {
	c.BaseController.Prepare()
	c.checkLogin()

}
func (c *MoldAnalyzeController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "moldanalyze/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "moldanalyze/index_footerjs.html"
}
