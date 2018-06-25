package validator

import (
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego"
	"bitsync/controllers"
)

type BaseValidator struct {
}

// 验证数据
func (v *BaseValidator) Validate(c *controllers.BaseController, o interface{}) {
	valid := validation.Validation{}
	ok, err := valid.Valid(o)

	if err != nil {
		beego.Error(err)
		c.OutputDefined(500, "", "未知错误")
	}

	if !ok {
		for _, err := range valid.Errors {
			msg := err.Key + " -> " + err.Message
			c.OutputDefined(400, "", msg)
			break
		}
	}
}
