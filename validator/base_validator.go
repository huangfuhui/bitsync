package validator

import (
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego"
	"bitsync/controllers"
)

type BaseValidator struct {
}

func init() {
	// 初始化错误信息
	validation.SetDefaultMessage(map[string]string{
		"Required":     "不能为空",
		"Min":          "最小值是%d",
		"Max":          "最大值是%d",
		"Range":        "范围是从%d到%d",
		"MinSize":      "最小长度是%d",
		"MaxSize":      "最大长度是%d",
		"Length":       "长度必须是%d",
		"Alpha":        "必须是字符",
		"Numeric":      "必须是数字",
		"AlphaNumeric": "必须是字符或数字",
		"Match":        "必须匹配正则%s",
		"NoMatch":      "必须不匹配%s",
		"AlphaDash":    "必须是字符或数字或横杠或下划线",
		"Email":        "必须是邮箱地址",
		"IP":           "必须是IP地址",
		"Base64":       "必须是base64编码",
		"Mobile":       "必须是手机号码",
		"Tel":          "必须是座机号码",
		"Phone":        "必须是座机或者手机号码",
		"ZipCode":      "必须是邮政编码",
	})
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
			msg := err.Field + err.Message
			c.OutputDefined(400, "", msg)
			break
		}
	}
}
