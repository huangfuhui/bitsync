package member

import (
	"bitsync/models"
	"bitsync/object/member"
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type MemberModel struct {
	models.BaseModel
}

// 新建会员
func (m *MemberModel) NewMember(uid int, name, handset, email string, sex int, avatarUrl string, BirthDay time.Time) (member.Member, error) {
	newMember := member.Member{
		UID:       uid,
		Name:      name,
		Handset:   handset,
		Email:     email,
		Sex:       sex,
		AvatarUrl: avatarUrl,
		Birthday:  BirthDay,
	}

	id, err := orm.NewOrm().Insert(&newMember)
	newMember.Id = int(id)

	return newMember, err
}

// 获取会员信息
func (m *MemberModel) Get(UID int) (member.Member, error) {
	memberInfo := member.Member{}
	memberInfo.UID = UID
	err := orm.NewOrm().Read(&memberInfo, "uid")

	return memberInfo, err
}

// 更新会员信息
func (m *MemberModel) Update(UID int, params orm.Params) error {
	_, err := orm.NewOrm().QueryTable("member").Filter("uid", strconv.FormatInt(int64(UID), 10)).Update(params)
	return err
}
