package wechat

import (
	"bitsync/controllers"
	"sort"
	"crypto/sha1"
)

type IndexController struct {
	controllers.BaseController
}

// 验证和接入微信开发平台
func (c *IndexController) Auth() {
	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")

	// 将参数排序和拼接
	str := sort.StringSlice{signature, timestamp, nonce, echostr}
	sort.Sort(str)
	sortStr := ""
	for _, v := range str {
		sortStr += v
	}

	// 进行sha1加密
	sh := sha1.New()
	sh.Write([]byte(sortStr))
	encryptStr := sh.Sum(nil)

	// 将本地计算的签名和微信传递过来的签名进行对比
	if string(encryptStr) == signature {
		c.Ctx.WriteString(signature)
	}

	c.Ctx.WriteString("Invalid signature")
}

// 关注时的欢迎语
func (c *IndexController) Welcome() {

}
