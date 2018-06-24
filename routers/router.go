package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
	"github.com/astaxie/beego/context"
	"bitsync/controllers/account"
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

	// 账号管理
	accountNs := beego.NewNamespace("/account",
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),

		beego.NSRouter("/register", &account.AccountController{}, "get:Register"),
	)

	beego.AddNamespace(wechatNs, smsNs, accountNs)
}
