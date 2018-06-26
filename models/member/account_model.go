package member

import (
	"bitsync/models"
	"bitsync/object/member"
	"time"
	"github.com/astaxie/beego/orm"
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

	query := `
update account 
set UID = 10000 + id
where id = ?
and account = ?
`
	_, err = o.Raw(query, newAccount.Id, account).Exec()
	if err != nil {
		o.Rollback()

		return 0, err
	}

	o.Commit()

	return newAccount.UID, nil
}

// 更新登录信息
func (m *AccountModel) Login(id int, loginTime, loginIp string) error {
	query := `
update account 
set last_login_time = login_time, last_login_ip = login_ip, login_time = ?, login_ip = ?
where id = ?
`
	_, err := orm.NewOrm().Raw(query, loginTime, loginIp, id).Exec()

	return err
}

// 验证账号密码
func (m *AccountModel) Verify(account, password string) bool {
	a := new(member.Account)
	a.Account = account
	a.Password = password

	err := orm.NewOrm().Read(a, "Account", "Password")
	if err == nil && a.UID > 0 {
		return true
	}

	return false
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
