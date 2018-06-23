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
func (model *MemberModel) NewMember(accountId int, name, handset, email string, sex int, avatarUrl string, BirthDay time.Time) (member.Member, error) {
	newMember := member.Member{
		AccountId: accountId,
		Name:      name,
		HandSet:   handset,
		Email:     email,
		Sex:       sex,
		AvatarUrl: avatarUrl,
		BirthDay:  BirthDay,
	}

	id, err := orm.NewOrm().Insert(&newMember)
	newMember.Id = int(id)

	return newMember, err
}

// 获取会员信息
func (model *MemberModel) Get(id int) (member.Member, error) {
	memberInfo := member.Member{}
	memberInfo.Id = id
	err := orm.NewOrm().Read(&memberInfo)

	return memberInfo, err
}
