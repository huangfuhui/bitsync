package Member

import (
	"bitsync/logic"
	"bitsync/models/member"
	"github.com/astaxie/beego"
	"strconv"
	"github.com/astaxie/beego/orm"
	"time"
)

type MemberLogic struct {
	logic.BaseLogic
}

// 查询用户信息
func (l *MemberLogic) Get() map[string]string {
	UID := l.GetUID()

	res := make(map[string]string, 10)

	m := member.MemberModel{}
	info, err := m.Get(UID)
	if err != nil || info.Id == 0 {
		beego.Error(err)
		l.Warn("查询用户信息失败")
		return res
	}

	res["uid"] = strconv.FormatInt(int64(info.UID), 10)
	res["name"] = info.Name
	res["handset"] = info.Handset
	res["email"] = info.Email
	res["sex"] = strconv.FormatInt(int64(info.Sex), 10)
	res["birthday"] = info.Birthday.Format("2006-01-02")
	res["avatar_url"] = info.AvatarUrl

	return res
}

// 更新用户信息
func (l *MemberLogic) Update(name, email string, sex int, birthday string) {
	UID := l.GetUID()

	birth, err := time.Parse("2006-01-02 15:04:05", birthday)
	if err != nil {
		beego.Error(err)
		l.Warn("更新失败")
		return
	}

	m := member.MemberModel{}
	err = m.Update(UID, orm.Params{
		"Name":     name,
		"Email":    email,
		"Sex":      sex,
		"Birthday": birth,
	})
	if err != nil {
		beego.Error(err)
		l.Warn("更新失败")
		return
	}
}
