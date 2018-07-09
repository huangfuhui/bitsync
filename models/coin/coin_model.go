package coin

import (
	"bitsync/models"
	"bitsync/object/coin"
	"github.com/astaxie/beego/orm"
)

type CoinModel struct {
	models.BaseModel
}

// 根据名称查询货币信息
func (m *CoinModel) GetByName(name string) (coin.Coin, error) {
	coinInfo := coin.Coin{Name: name}

	err := orm.NewOrm().Read(&coinInfo, "Name")
	if err != nil {
		return coin.Coin{}, err
	}

	return coinInfo, nil
}
