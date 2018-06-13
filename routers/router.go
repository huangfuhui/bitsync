package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
)

func init() {
	wechatNs := beego.NewNamespace("/wechat",
		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth;post:Dispatch"),
	)

	smsNs := beego.NewNamespace("/sms",
		beego.NSRouter("/index", &sms.IndexController{}, "get:Index"),
	)

	beego.AddNamespace(wechatNs, smsNs)
}
