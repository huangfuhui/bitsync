package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type response struct {
	code     int
	response interface{}
	msg      string
}

type BaseController struct {
	beego.Controller
}

// 正常响应输出JSON数据
func (c *BaseController) output(data interface{}) {
	responseData := new(response)
	responseData.code = http.StatusOK
	responseData.response = data
	responseData.msg = ""

	c.Data["json"] = responseData
	c.ServeJSONP()
}
