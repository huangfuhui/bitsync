package pay

import (
	"bitsync/models"
	"github.com/astaxie/beego/orm"
	"bitsync/object/pay"
)

type ComboModel struct {
	models.BaseModel
}

// 查询可用套餐
func (m *ComboModel) Get() ([]pay.Combo, error) {
	var combos []pay.Combo
	_, err := orm.NewOrm().
		QueryTable("combo").
		Filter("status", pay.COMBO_STAUTS_YES).
		All(&combos)
	if err != nil {
		return []pay.Combo{}, err
	}

	return combos, nil
}
