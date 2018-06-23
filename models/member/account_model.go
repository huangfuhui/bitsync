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
func (model *AccountModel) NewAccount(UID int, account, password, wechatOpeonid string) (int64, error) {
	newAccount := new(member.Account)
	newAccount.UID = UID
	newAccount.Account = account
	newAccount.Password = password
	newAccount.Status = member.STATUS_YES
	newAccount.WechatOpeonid = wechatOpeonid
	newAccount.RegisterTime = time.Now()

	return orm.NewOrm().Insert(&newAccount)
}

// 更新登录信息
func (model *AccountModel) Login(id int, loginTime, loginIp string) error {
	query := `
update account 
set last_login_time = login_time, last_login_ip = login_ip, login_time = ?, login_ip = ?
where id = ?
`
	_, err := orm.NewOrm().Raw(query, loginTime, loginIp, id).Exec()

	return err
}

// 验证账号密码
func (model *AccountModel) Verify(UID int, account, password string) bool {
	acc := new(member.Account)
	acc.UID = UID

	err := orm.NewOrm().Read(&acc, "UID")
	if err == nil && acc.Account == account && acc.Password == password {
		return true
	}

	return false
}
