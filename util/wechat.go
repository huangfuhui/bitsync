package util

import (
	"github.com/astaxie/beego"
	"sort"
	"crypto/sha1"
	"fmt"
)

var token string

func init() {
	token = beego.AppConfig.String("wechat_token")
}

// 验证签名
func AuthVerify(signature, timestamp, nonce, echostr string) (string, bool) {
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
		return echostr, true
	}

	return "Invalid Signature.", false
}
