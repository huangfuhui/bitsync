package middleware

import (
	"bitsync/util"
	"strconv"
	"time"
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
	redis := util.Cli{}
	redis.Select(0)
	res, _ := redis.Get("token:" + token)
	if res != "" {
		redis.Set("token:"+token, strconv.FormatInt(time.Now().Unix(), 64))
		return true
	}
	return false
}
