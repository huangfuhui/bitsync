package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
	"encoding/base64"
	"strings"
	"strconv"
)

type response struct {
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
	Msg      string      `json:"msg"`
}

type BaseController struct {
	beego.Controller
}

// 获取当前登录用户的UID
func (c *BaseController) GetUID() (UID int) {
	token := c.Ctx.Request.Header.Get("token")
	if token == "" {
		return 0
	}

	// 解密token
	res, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		beego.Error(err)

		return 0
	} else {
		decodeToken := strings.Split(string(res), ":")
		uid, err := strconv.ParseInt(decodeToken[0], 10, 0)
		if err != nil {
			beego.Error(err)

			return 0
		}

		return int(uid)
	}

	return 0
}

// 获取当前登录用户的账号
func (c *BaseController) GetAccount() (handset string) {
	token := c.Ctx.Request.Header.Get("token")
	if token == "" {
		return ""
	}

	// 解密token
	res, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		beego.Error(err)

		return ""
	} else {
		decodeToken := strings.Split(string(res), ":")

		return decodeToken[1]
	}

	return ""
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

// 服务器错误
func (c *BaseController) ServerError() {
	responseData := new(response)
	responseData.Code = http.StatusInternalServerError
	responseData.Response = ""
	responseData.Msg = "服务器错误"

	c.Ctx.Output.Status = http.StatusInternalServerError
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
