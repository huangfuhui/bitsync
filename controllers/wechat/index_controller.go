package wechat

import (
	"bitsync/controllers"
	"bitsync/util"
	"github.com/astaxie/beego"
)

type IndexController struct {
	controllers.BaseController
}

// 验证和接入微信开发平台
func (c *IndexController) Auth() {
	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")

	res, ok := util.AuthVerify(signature, timestamp, nonce, echostr)

	if ok {
		beego.Info("【微信】签名验证通过,成功接入微信开发平台.")
	} else {
		beego.Info("【微信】签名验证不通过,拒绝接入.")
	}

	c.Ctx.WriteString(res)
}

// 分发事件
func (c *IndexController) Dispatch() {
	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")

	res, ok := util.AuthVerify(signature, timestamp, nonce, echostr)

	if !ok {
		c.Ctx.WriteString(res)

		return
	}

	xmlBody := c.Ctx.Input.RequestBody

	// 解析XML
	base, err := util.ParseBase(xmlBody)

	if err != nil {
		beego.Error("【微信】解析推送数据失败: " + "\n" + string(xmlBody))
		beego.Error(err)

		c.Ctx.WriteString("")

		return
	} else {
		beego.Debug("【微信】解析推送数据成功, openid:" + base.FromUserName + ", msgtype:" + base.MsgType)
	}

	var response string
	switch base.MsgType {
	case util.MSGTYPE_TEXT:
		response = c.autoReply(xmlBody)
	case util.MSGTYPE_EVENT:
		response = c.event(xmlBody)
	default:
		response = "success"
	}

	c.Ctx.WriteString(response)
}

// 事件推送
func (c *IndexController) event(xmlBody []byte) string {
	res, err := util.ParseSubEvent(xmlBody)

	if err != nil {
		beego.Error("【微信】解析事件推送失败: " + "\n" + string(xmlBody))
		beego.Error(err)

		return ""
	}

	// TODO:消息排重

	// 拼接回复信息
	replay, err := util.ReplayTextMsg(res.BaseData.FromUserName, "终于等到你，还好我没放弃~~~")
	if err != nil {
		return ""
	}

	return string(replay)
}

// 自动回复
func (c *IndexController) autoReply(xmlBody []byte) string {
	res, err := util.ParseMsg(xmlBody)

	if err != nil {
		beego.Error("【微信】解析用户文本消息失败: " + "\n" + string(xmlBody))
		beego.Error(err)

		return ""
	}

	// TODO:消息排重

	// 拼接回复信息
	replay, err := util.ReplayTextMsg(res.BaseData.FromUserName, "收到的测试信息: "+res.Content)
	if err != nil {
		return ""
	}

	return string(replay)
}
