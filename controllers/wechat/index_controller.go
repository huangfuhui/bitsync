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
		beego.Info("微信签名验证通过,成功接入微信开发平台.")
	} else {
		beego.Info("微信签名验证不通过,拒绝接入.")
	}

	c.Ctx.WriteString(res)
}

// 关注时的欢迎语
func (c *IndexController) Welcome() {

}

// 自动回复
func (c *IndexController) AutoReply() {
	xmlBody := c.Ctx.Input.RequestBody

	res, err := util.ParseMsg(xmlBody)

	if err != nil {
		beego.Error("解析用户文本消息失败: " + "\n" + string(xmlBody))
		beego.Error(err)

		c.Ctx.WriteString("")

		return
	}

	beego.Debug("解析用户文本消息成功, openid:" + res.FromUserName + ", content:" + res.Content)

	// TODO:消息排重

	replay, err := util.ReplayTextMsg(res.FromUserName, "收到的测试信息: "+res.Content)
	if err != nil {
		c.Ctx.WriteString("")

		return
	}

	replayMsg := string(replay)
	c.Ctx.WriteString(replayMsg)
}
