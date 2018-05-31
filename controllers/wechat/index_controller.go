package wechat

import (
	"bitsync/controllers"
	"sort"
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego"
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

	token := beego.AppConfig.String("wechat_token")

	// 将参数排序和拼接
	str := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(str)
	sortStr := ""
	for _, v := range str {
		sortStr += v
	}

	// 进行sha1加密
	sh := sha1.New()
	sh.Write([]byte(sortStr))
	encryptStr := fmt.Sprintf("%x", sh.Sum(nil))

	// 将本地计算的签名和微信传递过来的签名进行对比
	if encryptStr == signature {
		c.Ctx.WriteString(echostr)

		return
	}

	c.Ctx.WriteString("Invalid Signature")
}

// 关注时的欢迎语
func (c *IndexController) Welcome() {

}
