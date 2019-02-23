package controllers

type MachineAnalyzeController struct {
	BaseController
}

func (c *MachineAnalyzeController) Prepare() {
	c.BaseController.Prepare()
	c.checkLogin()

}
func (c *MachineAnalyzeController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "machineanalyze/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "machineanalyze/index_footerjs.html"
}
