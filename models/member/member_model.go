package member

import (
	"bitsync/models"
	"bitsync/object/member"
	"time"
	"github.com/astaxie/beego/orm"
)

type MemberModel struct {
	models.BaseModel
}

// 新建会员
func (model *MemberModel) NewMember(uid int, name, handset, email string, sex int, avatarUrl string, BirthDay time.Time) (member.Member, error) {
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
func (model *MemberModel) Get(UID int) (member.Member, error) {
	memberInfo := member.Member{}
	memberInfo.UID = UID
	err := orm.NewOrm().Read(&memberInfo, "uid")

	return memberInfo, err
}
