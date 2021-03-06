package validator

import (
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego"
	"bitsync/controllers"
	"time"
	"strings"
	"regexp"
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
		"Email":        "必须是有效的邮箱地址",
		"IP":           "必须是IP地址",
		"Base64":       "必须是base64编码",
		"Mobile":       "必须是有效的手机号码",
		"Tel":          "必须是有效的座机号码",
		"Phone":        "必须是有效的座机或者手机号码",
		"ZipCode":      "必须是有效的邮政编码",
	})

	// 日期与时间校验
	validation.AddCustomFunc("Date", func(v *validation.Validation, obj interface{}, key string) {
		if obj.(string) == "" {
			return
		}

		_, err := time.Parse("2006-01-02", obj.(string))
		if err != nil {
			e := validation.Error{
				Key:     key,
				Name:    "Date",
				Field:   strings.Split(key, ".")[0],
				Value:   obj.(string),
				Message: "必须是 2006-01-02 格式的时间字符串",
			}
			v.Errors = append(v.Errors, &e)
		}
	})
	validation.AddCustomFunc("Time", func(v *validation.Validation, obj interface{}, key string) {
		if obj.(string) == "" {
			return
		}

		_, err := time.Parse("15:04:05", obj.(string))
		if err != nil {
			e := validation.Error{
				Key:     key,
				Name:    "Time",
				Field:   strings.Split(key, ".")[0],
				Value:   obj.(string),
				Message: "必须是 15:04:05 格式的时间字符串",
			}
			v.Errors = append(v.Errors, &e)
		}
	})
	validation.AddCustomFunc("DateTime", func(v *validation.Validation, obj interface{}, key string) {
		if obj.(string) == "" {
			return
		}

		_, err := time.Parse("2006-01-02 15:04:05", obj.(string))
		if err != nil {
			e := validation.Error{
				Key:     key,
				Name:    "DateTime",
				Field:   strings.Split(key, ".")[0],
				Value:   obj.(string),
				Message: "必须是 2006-01-02 15:04:05 格式的时间字符串",
			}
			v.Errors = append(v.Errors, &e)
		}
	})

	// 邮件校验
	validation.AddCustomFunc("Mail", func(v *validation.Validation, obj interface{}, key string) {
		value := obj.(string)
		if value == "" {
			return
		}

		ok, err := regexp.MatchString(`^[\w!#$%&'*+/=?^_`+"`"+`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`+"`"+`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`, value)
		if !ok || err != nil {
			e := validation.Error{
				Key:     key,
				Name:    "Mail",
				Field:   strings.Split(key, ".")[0],
				Value:   value,
				Message: "必须是有效邮箱地址",
			}
			v.Errors = append(v.Errors, &e)
		}
	})

	// 数字校验
	validation.AddCustomFunc("Num", func(v *validation.Validation, obj interface{}, key string) {
		value := obj.(string)
		if value == "" {
			return
		}

		ok, err := regexp.MatchString(`^-?[1-9][0-9]*\.?[0-9]*$`, value)
		if !ok || err != nil {
			e := validation.Error{
				Key:     key,
				Name:    "Num",
				Field:   strings.Split(key, ".")[0],
				Value:   value,
				Message: "必须是数字",
			}
			v.Errors = append(v.Errors, &e)
		}
	})
}

// 验证数据
func (v *BaseValidator) Validate(c *controllers.BaseController, o interface{}) bool {
	valid := validation.Validation{}
	ok, err := valid.Valid(o)

	if err != nil {
		beego.Error(err)
		c.ServerError()
		return false
	}

	if !ok {
		for _, err := range valid.Errors {
			msg := err.Field + err.Message
			c.OutputDefined(400, "", msg)
			break
		}
		return false
	}

	return true
}
