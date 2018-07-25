package market

import (
	"bitsync/logic"
	"bitsync/models/coin"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ExchangeLogic struct {
	logic.BaseLogic
}

// 获取交易所列表
func (l *ExchangeLogic) List() []orm.Params {
	m := coin.ExchangeModel{}
	res, err := m.List()
	if err != nil {
		beego.Error(err)
		return []orm.Params{}
	}
	return res
}
