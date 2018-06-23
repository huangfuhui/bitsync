package validator

import (
	"github.com/astaxie/beego/validation"
	"bitsync/controllers"
)

type BaseValidator struct {
}

// 验证数据
func (v *BaseValidator) Validate(c *controllers.BaseController, o *interface{}) {
	valid := validation.Validation{}
	ok, _ := valid.Valid(&o)

	if !ok {
		err := valid.Errors[0]
		msg := err.Key + " -> " + err.Message
		c.OutPutDefined(401, []string{}, msg)
	}
}
