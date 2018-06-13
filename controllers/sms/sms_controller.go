package sms

import "bitsync/controllers"

type IndexController struct {
	controllers.BaseController
}

func (c *IndexController) Index() {
	c.Ctx.WriteString("hello")
}
