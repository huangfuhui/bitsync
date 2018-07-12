package routers

import (
	"github.com/astaxie/beego"
	"bitsync/controllers/wechat"
	"bitsync/controllers/sms"
	"github.com/astaxie/beego/context"
	"bitsync/controllers/account"
	"github.com/astaxie/beego/plugins/cors"
	"bitsync/controllers/pay"
	"bitsync/middleware"
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

	// 微信授权
	wechatNs := beego.NewNamespace("/wechat",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/auth", &wechat.IndexController{}, "get:Auth;post:Dispatch"),
	)

	smsNs := beego.NewNamespace("/sms",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/index", &sms.IndexController{}, "get:Index"),
	)

	// 账号管理
	loginNs := beego.NewNamespace("/account",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth()
		}),

		beego.NSRouter("/register", &account.AccountController{}, "post:Register"),
		beego.NSRouter("/registerPin", &account.AccountController{}, "post:RegisterPin"),
		beego.NSRouter("/login", &account.AccountController{}, "post:Login"),
	)

	// 账号管理
	accountNs := beego.NewNamespace("/account",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/modifyPassword", &account.AccountController{}, "post:ModifyPassword"),
		beego.NSRouter("/passwordPin", &account.AccountController{}, "post:PasswordPin"),
		beego.NSRouter("/resetPassword", &account.AccountController{}, "post:ResetPassword"),
	)

	// 用户管理
	memberNs := beego.NewNamespace("/member",
		beego.NSBefore(func(ctx *context.Context) {
			m := middleware.Base{ctx}
			m.Auth(middleware.Auth{})
		}),

		beego.NSRouter("/get", &account.MemberController{}, "get:Get"),
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
	)

	beego.AddNamespace(wechatNs, smsNs, loginNs, accountNs, memberNs, comboNs, taskNs)
}
