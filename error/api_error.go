package error

import "bitsync/controllers"

type ApiError struct {
}

// 参数错误
func (err *ApiError) ApiErr(c *controllers.BaseController, error string) {
	c.OutPutDefined(400, []string{}, error)
}
