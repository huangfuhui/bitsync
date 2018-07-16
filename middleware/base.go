package middleware

import (
	"github.com/astaxie/beego/context"
	"net/http"
	"bitsync/controllers"
)

type Base struct {
	*context.Context
}

// 执行中间件
func (base *Base) Auth(middlewares ...interface{}) {
	responseData := new(controllers.Response)
	for _, middleware := range middlewares {
		res := true
		switch middleware.(type) {
		case Auth:
			auth := Auth{base}
			res = auth.Verify()

			responseData.Code = -400
		}

		if !res {
			responseData.Response = ""
			responseData.Msg = "请求错误"

			base.Output.Status = http.StatusBadRequest
			base.Output.JSON(responseData, false, false)
			break
		}
	}
}
