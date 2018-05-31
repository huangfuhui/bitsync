package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
)

func init() {
	ns := beego.NewNamespace("/wechat",
		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth"),
		beego.NSRouter("/welcome", &wechat.IndexController{}, "get:Welcome"),
	)

	beego.AddNamespace(ns)
}
