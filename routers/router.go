package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
	"github.com/astaxie/beego/context"
	"bitsync/controllers/account"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	// 允许跨域请求
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

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

		beego.NSRouter("/register", &account.AccountController{}, "post:Register"),
		beego.NSRouter("/registerPin", &account.AccountController{}, "post:RegisterPin"),
		beego.NSRouter("/login", &account.AccountController{}, "post:Login"),
		beego.NSRouter("/modifyPassword", &account.AccountController{}, "post:ModifyPassword"),
		beego.NSRouter("/passwordPin", &account.AccountController{}, "post:PasswordPin"),
		beego.NSRouter("/resetPassword", &account.AccountController{}, "post:ResetPassword"),
	)

	// 用户管理
	memberNs := beego.NewNamespace("/member",
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),

		beego.NSRouter("/get", &account.MemberController{}, "get:Get"),
	)

	beego.AddNamespace(wechatNs, smsNs, accountNs, memberNs)
}
