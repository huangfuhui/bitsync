package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"bitsync/controllers/pay"
	"bitsync/middleware"
	"bitsync/controllers/member"
	"bitsync/controllers/market"
)

func init() {
	// 允许跨域请求
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Token", "Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	// 微信授权
	wechatNs := beego.NewNamespace("/wechat",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth;post:Dispatch"),
	)

	// 注册登录
	loginNs := beego.NewNamespace("/account",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/register", &member.AccountController{}, "post:Register"),
		beego.NSRouter("/registerPin", &member.AccountController{}, "post:RegisterPin"),
		beego.NSRouter("/login", &member.AccountController{}, "post:Login"),
		beego.NSRouter("/passwordPin", &member.AccountController{}, "post:PasswordPin"),
		beego.NSRouter("/resetPassword", &member.AccountController{}, "post:ResetPassword"),
	)

	// 账号管理
	accountNs := beego.NewNamespace("/accountManage",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/logout", &member.AccountController{}, "get:Logout"),
		beego.NSRouter("/modifyPassword", &member.AccountController{}, "post:ModifyPassword"),
	)

	// 用户管理
	memberNs := beego.NewNamespace("/member",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/get", &member.MemberController{}, "get:Get"),
		beego.NSRouter("/update", &member.MemberController{}, "post:Update"),
	)

	// 短信
	smsNs := beego.NewNamespace("/sms",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/wallet", &sms.IndexController{}, "get:Wallet"),
	)

	// 短信套餐
	comboNs := beego.NewNamespace("/combo",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/get", &pay.ComboController{}, "get:Get"),
	)

	// 预警任务
	taskNs := beego.NewNamespace("/task",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/add", &sms.TaskController{}, "post:Add"),
		beego.NSRouter("/cancel", &sms.TaskController{}, "post:Cancel"),
		beego.NSRouter("/list", &sms.TaskController{}, "post:List"),
		beego.NSRouter("/get", &sms.TaskController{}, "post:Get"),
	)

	// 市场行情
	marketNs := beego.NewNamespace("/market",
		beego.NSBefore(func(ctx *context.Context) {
		}),

		beego.NSRouter("/exchange", &market.ExchangeController{}, "get:List"),
	)

	beego.AddNamespace(wechatNs, smsNs, loginNs, accountNs, memberNs, comboNs, taskNs, marketNs)
}
