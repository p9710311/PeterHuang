package controllers

type ScheduleAnalyzeController struct {
	BaseController
}

func (c *ScheduleAnalyzeController) Prepare() {
	c.BaseController.Prepare()
	c.checkLogin()

}
func (c *ScheduleAnalyzeController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	// c.LayoutSections["headcssjs"] = "scheduleanalyze/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "scheduleanalyze/index_footerjs.html"
}
