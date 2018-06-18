package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
	"github.com/astaxie/beego/context"
	"bitsync/middleware"
)

func init() {
	wechatNs := beego.NewNamespace("/wechat",
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),

		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth;post:Dispatch"),
	)

	smsNs := beego.NewNamespace("/sms",
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),

		beego.NSRouter("/index", &sms.IndexController{}, "get:Index"),
	)

	indexNs := beego.NewNamespace("/index",
		beego.NSCond(func(ctx *context.Context) bool {
			base := middleware.Base{Context: ctx}
			return base.Auth(
				middleware.Auth{},
			)
		}),

	)

	beego.AddNamespace(wechatNs, smsNs, indexNs)
}
