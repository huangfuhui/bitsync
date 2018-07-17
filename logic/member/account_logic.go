package member

import (
	"bitsync/logic"
	"regexp"
	"bitsync/services"
	"bitsync/models/member"
	"github.com/astaxie/beego"
	"crypto/md5"
	"io"
	"fmt"
	"bitsync/util"
	"strconv"
	"time"
	"strings"
	"encoding/base64"
)

type AccountLogic struct {
	logic.BaseLogic
}

// 注册
func (l *AccountLogic) Register(handset, password, pin string) (res map[string]string) {
	account := member.AccountModel{}
	exists := account.Exists(handset)
	if exists > 0 {
		l.Warn("账号已经存在")
		return
	}

	// 校验密码强度
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", password)
	if !match {
		l.BadRequest("密码长度必须是8-15个字符")
		return
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个字母")
		return
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个数字")
		return
	}

	// 校验验证码
	sms := services.PinService{}
	_, err := sms.Validate(services.PIN_REGISTER, handset, pin)
	if err != nil {
		l.Warn(err.Error())
		return
	}

	// 密码加密
	salt := beego.AppConfig.String("salt")
	w := md5.New()
	io.WriteString(w, salt+password)
	password = fmt.Sprintf("%x", w.Sum(nil))

	// 注册
	var UID int
	UID, err = account.NewAccount(handset, password, "")
	if err != nil {
		beego.Error(err)
		l.ServerError()
		return
	}

	// 更新登录信息
	remoteAddr := strings.Split(l.Ctx.Request.RemoteAddr, ":")
	err = account.Login(UID, time.Now().Format("2006-01-02 15:04:05"), remoteAddr[0])
	if err != nil {
		beego.Error(err)
	}

	random := util.Random{}
	randomNum := random.Rand(100000, 999999)

	// 生成token
	tokenKey := strconv.FormatInt(int64(UID), 10) + ":" + handset + ":" + strconv.FormatInt(randomNum, 10)
	token := base64.StdEncoding.EncodeToString([]byte(tokenKey))

	// 保存token
	db, _ := beego.AppConfig.Int("redis_db_token")
	redis := util.Cli{}
	redis.Select(db)
	key := "token:" + handset
	redis.Set(key, token)
	redis.SetEx(key, 3600)

	return map[string]string{"uid": strconv.FormatInt(int64(UID), 10), "token": token}
}

// 发送注册验证码
func (l *AccountLogic) RegisterPin(handset string) {
	sms := services.PinService{}
	_, err := sms.Send(services.PIN_REGISTER, handset)
	if err != nil {
		l.Warn(err.Error())
		return
	}
}

// 登录
func (l *AccountLogic) Login(handset, password string) (res map[string]string) {
	salt := beego.AppConfig.String("salt")
	w := md5.New()
	io.WriteString(w, salt+password)
	password = fmt.Sprintf("%x", w.Sum(nil))

	// 验证账号密码
	account := member.AccountModel{}
	UID := account.Verify(handset, password)
	if UID == 0 {
		l.BadRequest("账号或密码不正确")
		return
	}

	random := util.Random{}
	randomNum := random.Rand(100000, 999999)

	// 生成token
	tokenKey := strconv.FormatInt(int64(UID), 10) + ":" + handset + ":" + strconv.FormatInt(randomNum, 10)
	token := base64.StdEncoding.EncodeToString([]byte(tokenKey))

	// 保存token
	db, _ := beego.AppConfig.Int("redis_db_token")
	redis := util.Cli{}
	redis.Select(db)
	key := "token:" + handset
	redis.Set(key, token)
	redis.SetEx(key, 3600)

	return map[string]string{"token": token}
}

// 修改密码
func (l *AccountLogic) ModifyPassword(oldPwd, newPwd string) {
	UID := l.GetUID()
	handset := l.GetAccount()
	if UID == 0 || handset == "" {
		l.Warn("修改密码失败")

		return
	}

	salt := beego.AppConfig.String("salt")
	w := md5.New()
	io.WriteString(w, salt+oldPwd)
	oldPwd = fmt.Sprintf("%x", w.Sum(nil))

	// 验证账号旧密码
	account := member.AccountModel{}
	exists := account.Verify(handset, oldPwd)
	if exists == 0 {
		l.BadRequest("旧密码不正确")
		return
	}

	// 校验新密码强度
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", newPwd)
	if !match {
		l.BadRequest("密码长度必须是8-15个字符")
		return
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", newPwd)
	if !match {
		l.BadRequest("密码至少包含一个字母")
		return
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", newPwd)
	if !match {
		l.BadRequest("密码至少包含一个数字")
		return
	}

	// 新密码加密
	w = md5.New()
	io.WriteString(w, salt+newPwd)
	newPwd = fmt.Sprintf("%x", w.Sum(nil))

	// 修改密码
	err := account.ModifyPassword(UID, oldPwd, newPwd)
	if err != nil {
		beego.Error(err)

		l.Warn("修改密码失败")
	}
}

// 发送重置密码验证码
func (l *AccountLogic) PasswordPin(handset string) {
	sms := services.PinService{}
	_, err := sms.Send(services.PIN_RESET_PASSWORD, handset)
	if err != nil {
		l.Warn(err.Error())
		return
	}
}

// 重置密码
func (l *AccountLogic) ResetPassword(handset, pin, newPwd string) {
	// 校验验证码
	sms := services.PinService{}
	_, err := sms.Validate(services.PIN_RESET_PASSWORD, handset, pin)
	if err != nil {
		l.Warn(err.Error())
		return
	}

	// 校验新密码强度
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", newPwd)
	if !match {
		l.BadRequest("密码长度必须是8-15个字符")
		return
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", newPwd)
	if !match {
		l.BadRequest("密码至少包含一个字母")
		return
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", newPwd)
	if !match {
		l.BadRequest("密码至少包含一个数字")
		return
	}

	// 新密码加密
	w := md5.New()
	salt := beego.AppConfig.String("salt")
	io.WriteString(w, salt+newPwd)
	newPwd = fmt.Sprintf("%x", w.Sum(nil))

	// 修改密码
	account := member.AccountModel{}
	err = account.ResetPassword(handset, newPwd)
	if err != nil {
		beego.Error(err)

		l.Warn("重置密码失败")
		return
	}
}

// 退出登录
func (l *AccountLogic) Logout() {
	handset := l.GetAccount()

	db, _ := beego.AppConfig.Int("redis_db_token")
	redis := util.Cli{}
	redis.Select(db)
	key := "token:" + handset
	err := redis.Del(key)
	if err != nil {
		beego.Error(err)
		l.Warn("退出登录失败")
		return
	}
}
