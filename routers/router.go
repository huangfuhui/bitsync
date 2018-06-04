package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
)

func init() {
	ns := beego.NewNamespace("/wechat",
		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth;post:Dispatch"),
	)

	beego.AddNamespace(ns)
}
