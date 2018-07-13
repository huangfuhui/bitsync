package middleware

import (
	"bitsync/util"
	"encoding/base64"
	"github.com/astaxie/beego"
	"strings"
)

type Auth struct {
	*Base
}

// 校验token有效性
func (auth *Auth) Verify() bool {
	token := auth.Request.Header.Get("token")
	if token == "" {
		return false
	}

	// 解密token
	res, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	} else {
		decodeToken := strings.Split(string(res), ":")

		db, _ := beego.AppConfig.Int("redis_db_token")
		redis := util.Cli{}
		redis.Select(db)
		localToken, err := redis.Get("token:" + decodeToken[1])
		if err != nil {
			return false
		} else if localToken != token {
			return false
		}

		// 刷新token有效期
		redis.SetEx("token:"+decodeToken[1], 3600)

		return true
	}

	return false
}
