package member

import (
	"bitsync/models"
	"bitsync/object/member"
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type AccountModel struct {
	models.BaseModel
}

// 新增账号
func (m *AccountModel) NewAccount(account, password, wechatOpeonid string) (UID int, err error) {
	o := orm.NewOrm()
	err = o.Begin()

	if err != nil {
		return 0, err
	}

	newAccount := new(member.Account)
	newAccount.Account = account
	newAccount.Password = password
	newAccount.Status = member.STATUS_YES
	newAccount.WechatOpenid = wechatOpeonid
	newAccount.RegisterTime = time.Now()

	_, err = o.Insert(newAccount)
	if err != nil {
		o.Rollback()

		return 0, err
	}

	newAccount.UID = 10000 + newAccount.Id

	query := `
update account 
set UID = ?
where id = ?
and account = ?
`
	_, err = o.Raw(query, newAccount.UID, newAccount.Id, account).Exec()
	if err != nil {
		o.Rollback()

		return 0, err
	}

	name := "bs" + strconv.FormatInt(int64(newAccount.UID), 10)
	newMember := member.Member{
		UID:       newAccount.UID,
		Name:      name,
		Handset:   account,
		Email:     "",
		Sex:       member.MEMBER_SEX_MALE,
		AvatarUrl: "",
		Birthday:  time.Now(),
	}

	_, err = o.Insert(&newMember)
	if err != nil {
		o.Rollback()

		return 0, err
	}

	o.Commit()

	return newAccount.UID, nil
}

// 更新登录信息
func (m *AccountModel) Login(UID int, loginTime, loginIp string) error {
	query := `
update account 
set last_login_time = login_time, last_login_ip = login_ip, login_time = ?, login_ip = ?
where uid = ?
`
	_, err := orm.NewOrm().Raw(query, loginTime, loginIp, UID).Exec()

	return err
}

// 验证账号密码
func (m *AccountModel) Verify(account, password string) (UID int) {
	a := new(member.Account)
	a.Account = account
	a.Password = password

	err := orm.NewOrm().Read(a, "Account", "Password")
	if err == nil && a.UID > 0 {
		return a.UID
	}

	return 0
}

// 判断账号是否存在
func (m *AccountModel) Exists(account string) (UID int) {
	a := new(member.Account)
	a.Account = account

	err := orm.NewOrm().Read(a, "Account")
	if err == nil && a.UID > 0 {
		return a.UID
	}

	return 0
}

// 修改密码
func (m *AccountModel) ModifyPassword(UID int, oldPwd, newPwd string) error {
	a := new(member.Account)
	a.UID = UID
	a.Password = oldPwd

	err := orm.NewOrm().Read(a, "UID", "Password")
	if err == nil {
		a.Password = newPwd
		_, err = orm.NewOrm().Update(a, "Password")
	}

	return err
}

// 重置密码
func (m *AccountModel) ResetPassword(handset, newPwd string) error {
	a := new(member.Account)
	a.Account = handset

	err := orm.NewOrm().Read(a, "Account")
	if err == nil {
		a.Password = newPwd
		_, err = orm.NewOrm().Update(a, "Password")
	}

	return err
}
