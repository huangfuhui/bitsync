package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type response struct {
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
	Msg      string      `json:"msg"`
}

type BaseController struct {
	beego.Controller
}

// 正常响应输出JSON数据
func (c *BaseController) Output(data interface{}) {
	responseData := new(response)
	responseData.Code = http.StatusOK
	responseData.Response = data
	responseData.Msg = ""

	c.Ctx.Output.Status = http.StatusOK
	c.Data["json"] = &responseData
	c.ServeJSON()
}

// 自定义响应输出JSON数据
func (c *BaseController) OutputDefined(code int, data interface{}, msg string) {
	responseData := new(response)
	responseData.Code = code
	responseData.Response = data
	responseData.Msg = msg

	c.Ctx.Output.Status = code
	c.Data["json"] = &responseData
	c.ServeJSON()
}

// 请求参数错误
func (c *BaseController) BadRequest(msg string) {
	responseData := new(response)
	responseData.Code = http.StatusBadRequest
	responseData.Response = ""
	responseData.Msg = msg

	c.Ctx.Output.Status = http.StatusBadRequest
	c.Data["json"] = &responseData
	c.ServeJSON()
}

// 业务错误
func (c *BaseController) Warn(msg string) {
	responseData := new(response)
	responseData.Code = -1
	responseData.Response = ""
	responseData.Msg = msg

	c.Ctx.Output.Status = http.StatusOK
	c.Data["json"] = &responseData
	c.ServeJSON()
}
