package middleware

import "github.com/astaxie/beego/context"

type Base struct {
	*context.Context
}

// 执行中间件
func (base *Base) Auth(middlewares ...interface{}) bool {
	for _, middleware := range middlewares {
		res := true
		switch middleware.(type) {
		case Auth:
			auth := Auth{base}
			res = auth.Verify()
		}

		if !res {
			return false
		}
	}

	return true
}
