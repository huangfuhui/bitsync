package coin

import (
	"bitsync/models"
	"github.com/astaxie/beego/orm"
)

type ExchangeModel struct {
	models.BaseModel
}

// 获取交易所列表
func (m *ExchangeModel) List() ([]orm.Params, error) {
	var exchanges []orm.Params
	_, err := orm.NewOrm().QueryTable("exchange").Values(&exchanges, "exchange_id", "name_cn", "name_en")
	if err != nil {
		return exchanges, err
	}

	return exchanges, nil
}
