package middleware

import "github.com/astaxie/beego/context"

type Base struct {
	Context *context.Context
}

// 执行中间件
func (base *Base) Auth(middlewares ...interface{}) bool {
	for _, middleware := range middlewares {
		res := true
		switch middleware := middleware.(type) {
		case Auth:
			res = middleware.Verify()
		}

		if !res {
			return false
		}
	}

	return false
}
